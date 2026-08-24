"use client"

import { useState } from "react"
import { createClient } from "@/lib/supabase/client"

import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"

type Profile = {
  username: string | null
  first_name: string | null
  last_name: string | null
  avatar_url: string | null
  zip_code: number | null
}

type SettingsFormProps = {
  userId: string
  profile: Profile | null
}

export function SettingsForm({
  userId,
  profile,
}: SettingsFormProps) {

  // ALL hooks go inside the component
  const [username, setUsername] = useState(profile?.username ?? "")
  const [firstName, setFirstName] = useState(profile?.first_name ?? "")
  const [lastName, setLastName] = useState(profile?.last_name ?? "")
  const [avatarUrl, setAvatarUrl] = useState(profile?.avatar_url ?? "")
  const [zipCode, setZipCode] = useState(
    profile?.zip_code?.toString() ?? ""
  )

  const [saving, setSaving] = useState(false)
  const [message, setMessage] = useState("")
  const [locating, setLocating] = useState(false)

  async function handleGetLocation() {
    setLocating(true)
    setMessage("")

    navigator.geolocation.getCurrentPosition(
      async (position) => {
        try {
          const { latitude, longitude } = position.coords

          const response = await fetch(
            `https://nominatim.openstreetmap.org/reverse?lat=${latitude}&lon=${longitude}&format=json`
          )

          if (!response.ok) {
            throw new Error("Failed to find ZIP code")
          }

          const data = await response.json()
          const postalCode = data.address?.postcode

          if (!postalCode) {
            throw new Error("ZIP code not found")
          }

          setZipCode(postalCode)
          setMessage("ZIP code found.")
        } catch (error) {
          console.error(error)
          setMessage("Could not determine your ZIP code.")
        } finally {
          setLocating(false)
        }
      },
      (error) => {
        console.error(error)
        setMessage("Location permission was denied.")
        setLocating(false)
      }
    )
  }

  async function handleSubmit(
    event: React.FormEvent<HTMLFormElement>
  ) {
    event.preventDefault()

    setSaving(true)
    setMessage("")

    const supabase = createClient()

    const { error } = await supabase
      .from("profiles")
      .update({
        username: username || null,
        first_name: firstName || null,
        last_name: lastName || null,
        avatar_url: avatarUrl || null,
        zip_code: zipCode ? Number(zipCode) : null,
      })
      .eq("id", userId)

    setSaving(false)

    if (error) {
      console.error(error)
      setMessage("Failed to save changes.")
      return
    }

    setMessage("Changes saved.")
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-6">

      <div className="space-y-2">
        <Label htmlFor="username">Username</Label>
        <Input
          id="username"
          value={username}
          onChange={(event) => setUsername(event.target.value)}
          placeholder="username"
        />
      </div>

      <div className="grid gap-6 sm:grid-cols-2">
        <div className="space-y-2">
          <Label htmlFor="first-name">First name</Label>
          <Input
            id="first-name"
            value={firstName}
            onChange={(event) => setFirstName(event.target.value)}
            placeholder="John"
          />
        </div>

        <div className="space-y-2">
          <Label htmlFor="last-name">Last name</Label>
          <Input
            id="last-name"
            value={lastName}
            onChange={(event) => setLastName(event.target.value)}
            placeholder="Doe"
          />
        </div>
      </div>

      <div className="space-y-2">
        <Label htmlFor="avatar-url">Avatar URL</Label>
        <Input
          id="avatar-url"
          type="url"
          value={avatarUrl}
          onChange={(event) => setAvatarUrl(event.target.value)}
          placeholder="https://..."
        />
      </div>

      <div className="space-y-2">
        <Label htmlFor="zip-code">ZIP code</Label>

        <div className="flex gap-2">
          <Input
            id="zip-code"
            type="text"
            inputMode="numeric"
            maxLength={5}
            value={zipCode}
            onChange={(event) => setZipCode(event.target.value)}
            placeholder="74103"
          />

          <Button
            type="button"
            variant="outline"
            onClick={handleGetLocation}
            disabled={locating}
          >
            {locating ? "Finding..." : "Use my location"}
          </Button>
        </div>
      </div>

      <div className="flex items-center gap-4">
        <Button type="submit" disabled={saving}>
          {saving ? "Saving..." : "Save changes"}
        </Button>

        {message && (
          <p className="text-sm text-muted-foreground">
            {message}
          </p>
        )}
      </div>

    </form>
  )
}
