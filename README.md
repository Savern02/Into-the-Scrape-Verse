# Cheap Chick

Be a Cheap Chick and save money on the groceries you buy every day.

## What is Cheap Chick

Cheap Chick is a grocery comparison site that gathers groceries from popular supermarkets. Either by individual items or providing a grocery list, Cheap Chick allows you to save money on your daily necessities.

## The Problem

Grocery prices vary significantly between stores, but comparing prices manually
is time-consuming. Existing grocery comparison tools often depend on static
datasets that become outdated. The UI is clunky and the data is usually inaccurate.

Cheap Chick automates the collection and validation of grocery prices so users
can compare current prices across supermarkets without checking each store
individually.

## Key features

* Self-repairing scrapers that are powered by Bright Data
* Go Service Pipeline for Data Validation and Ingest
* Agent-driven orchestration (run scrapers from a coding/automation agent)
* Modular TypeScript codebase for scraper logic, repair heuristics, and data processing powered by Supabase for Storage and Auth
* Easy to extend: add new site adapters, repair strategies, and downstream transforms

## Stack

* Language: TypeScript
* Next.js, Tailwind CSS, ShadCN, and Supabase.

## Project Directory

go-service ---> Contains Services for the Data Pipeline
web-app ---> Contains Next.js web app

## Quickstart (example)

Install dependencies:
npm install
Start the dev runner (example):
npm run dev
Build for production:
npm run build
Run tests:
npm test

## Configuration

Typical environment variables/config you’ll likely see or want:

## Architecture (high level)

* go agents/ — orchestrates scraping runs and repair attempts (agent loop detects failures and triggers repairs)
* scrapers/ — per-site adapters: navigation, selectors, and extraction rules
* processors/ — normalize and persist scraped data to storage
* auth/ contains all the Next.js logic

## Data Pipeline

The agent runs scraping sessions through site adapters; failures are surfaced to the repair component which applies a sequence of heuristics/tests to restore extraction, then the pipeline re-runs and persists the recovered output.

## Configuration

Create a `.env.local` file in `web-app`:

```env
NEXT_PUBLIC_SUPABASE_URL=...
NEXT_PUBLIC_SUPABASE_ANON_KEY=...
```

## Acknowledgements

Organized for the WeMakeDevs Scrape-Verse hackathon — see resources at https://www.wemakedevs.org/hackathons/scrape-verse/resourcesv

## License

TBD — add your chosen license (e.g., MIT) or check the repository’s LICENSE file.

