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

  let { data: profile, error: profileError } = await supabase
    .from("profiles")
    .select("username, first_name, last_name, avatar_url, zip_code")
    .eq("id", userId)
    .maybeSingle()

  if (profileError) {
    console.error(
      "Failed to fetch profile:",
      profileError.message
    )
  }

  // Create a blank profile if one doesn't exist
  if (!profile) {
    const { data: newProfile, error: insertError } =
      await supabase
        .from("profiles")
        .insert({
          id: userId,
        })
        .select(
          "username, first_name, last_name, avatar_url, zip_code"
        )
        .single()

    if (insertError) {
      console.error(
        "Failed to create profile:",
        insertError.message
      )
    } else {
      profile = newProfile
    }
  }

  return (
    <div className="max-w-2xl space-y-8">
      <div>
        <h1 className="text-3xl font-bold">
          Settings
        </h1>

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
