import { createClient } from "@/lib/supabase/server"
import { redirect } from "next/navigation"
import { SettingsForm } from "@/components/g-ui/settings/settings-form"

export default async function SettingsPage() {
  const supabase = await createClient()

  const { data, error } = await supabase.auth.getClaims()

  if (error || !data?.claims?.sub) {
    redirect("/auth/login")
  }

  const userId = data.claims.sub

  const { data: profile, error: profileError } = await supabase
    .from("profiles")
    .select("username, first_name, last_name, avatar_url, zip_code")
    .eq("id", userId)
    .single()

  if (profileError) {
    console.error("Failed to fetch profile:", profileError)
  }

  return (
    <div className="max-w-2xl space-y-8">
      <div>
        <h1 className="text-3xl font-bold">Settings</h1>
        <p className="text-muted-foreground">
          Manage your profile information.
        </p>
      </div>

      <SettingsForm
        userId={userId}
        profile={profile}
      />
    </div>
  )
}
