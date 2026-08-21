package brightdata

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

type Client struct {
	token string
	http  *http.Client
}

func New(token string) *Client {
	return &Client{token: token, http: &http.Client{Timeout: 2 * time.Minute}}
}

// Trigger fires a collector as a production endpoint. This is the line that
// makes best-practice 05 true for your project: the c_* ID is a live API, not
// a thing you click in a dashboard.
//
// Verify the exact path and payload shape against the API quickstart in the
// hackathon resources before you demo -- treat this as the shape, not gospel.
func (c *Client) Trigger(ctx context.Context, collectorID string, inputs []map[string]any) (string, error) {
	body, err := json.Marshal(inputs)
	if err != nil {
		return "", err
	}
	url := fmt.Sprintf("https://api.brightdata.com/dca/trigger?collector=%s&queue_next=1", collectorID)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 300 {
		return "", fmt.Errorf("trigger %s: %s: %s", collectorID, resp.Status, truncate(string(raw), 300))
	}

	var out struct {
		CollectionID string `json:"collection_id"`
		SnapshotID   string `json:"snapshot_id"`
	}
	_ = json.Unmarshal(raw, &out)
	if out.SnapshotID != "" {
		return out.SnapshotID, nil
	}
	return out.CollectionID, nil
}

// FetchSnapshot downloads the finished JSON for a run. The webhook usually
// hands you a URL directly; this is the fallback for polling.
func (c *Client) FetchSnapshot(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		raw, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("fetch snapshot: %s: %s", resp.Status, truncate(string(raw), 300))
	}
	return io.ReadAll(resp.Body)
}

// HealCommand builds the exact CLI invocation for a broken collector.
//
// The description matters more than anything else here. Bright Data's own docs
// are blunt about it: "Vague prompts produce vague heals." Max 1000 chars.
// Name what is wrong AND what correct output looks like.
func HealCommand(collectorID, whatBroke string) []string {
	return []string{
		"npx", "--yes", "--package", "@brightdata/cli",
		"brightdata", "scraper", "heal", collectorID, whatBroke,
		"--pretty",
	}
}

func ApproveCommand(collectorID string) []string {
	return []string{
		"npx", "--yes", "--package", "@brightdata/cli",
		"brightdata", "scraper", "approve", collectorID, "--pretty",
	}
}

// HealResult mirrors the CLI's JSON envelope.
type HealResult struct {
	CollectorID   string          `json:"collector_id"`
	Status        string          `json:"status"` // awaiting_approval | done | failed | rejected
	Prompt        string          `json:"prompt"`
	PreviewResult json.RawMessage `json:"preview_result"`
	DiffSummary   string          `json:"diff_summary"`
	ViewURL       string          `json:"view_url"`
	NextStep      string          `json:"next_step"`
	Error         string          `json:"error"`
	Raw           string          `json:"-"`
}

// Committed reports whether the fix is actually live. awaiting_approval is a
// success path, not a failure -- but the scraper is NOT fixed yet.
func (h HealResult) Committed() bool { return h.Status == "done" }

// Heal shells out to the CLI and stops at the approval gate by design.
//
// Self-healing is human-in-the-loop. A heal that returns awaiting_approval has
// a preview_result showing the rows the fixed scraper WOULD produce -- review
// those before committing. A failed heal is non-destructive: the existing
// scraper is unchanged and still works.
//
// Pass autoApprove=true to skip the gate entirely. Useful for the "scraper
// that fixes itself while you sleep" demo; risky everywhere else, because
// nobody looks at the preview.
func Heal(ctx context.Context, collectorID, whatBroke string, autoApprove bool, timeout time.Duration) (HealResult, error) {
	args := HealCommand(collectorID, whatBroke)
	if autoApprove {
		args = append(args, "--auto-approve")
	}
	return runCLI(ctx, args, timeout)
}

// Approve commits a fix left waiting by Heal. Pass reject=true to discard it
// and try again with a sharper prompt.
func Approve(ctx context.Context, collectorID string, reject bool, timeout time.Duration) (HealResult, error) {
	args := ApproveCommand(collectorID)
	if reject {
		args = append(args, "--reject")
	}
	return runCLI(ctx, args, timeout)
}

func runCLI(ctx context.Context, args []string, timeout time.Duration) (HealResult, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, args[0], args[1:]...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	runErr := cmd.Run()

	var res HealResult
	res.Raw = stdout.String() + stderr.String()
	// The envelope is written on every termination path, success or failure,
	// so try to parse it even when the process exited non-zero.
	if err := json.Unmarshal(stdout.Bytes(), &res); err != nil && runErr == nil {
		return res, fmt.Errorf("parse heal envelope: %w: %s", err, truncate(res.Raw, 300))
	}
	if runErr != nil {
		return res, fmt.Errorf("%s: %w: %s", args[4], runErr, truncate(res.Raw, 300))
	}
	return res, nil
}

// Describe turns null-rate stats into the sentence we hand to heal.
func Describe(nullRates map[string]float64, fieldSpec map[string]string) string {
	var broken []string
	for field, rate := range nullRates {
		if rate >= 0.5 {
			if desc, ok := fieldSpec[field]; ok && desc != "" {
				broken = append(broken, fmt.Sprintf("%q (%s)", field, desc))
			} else {
				broken = append(broken, fmt.Sprintf("%q", field))
			}
		}
	}
	if len(broken) == 0 {
		return "extraction quality dropped without a single field failing outright"
	}
	return fmt.Sprintf(
		"these fields came back empty on most rows after the site changed its markup: %s. "+
			"Rewrite the extraction against those descriptions.",
		strings.Join(broken, ", "))
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
