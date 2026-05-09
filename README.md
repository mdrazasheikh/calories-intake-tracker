# Calories Intake Tracker (Supabase-native deploy)

This project is now structured to run directly on **Supabase**:
- Frontend can be deployed from your Git repo as a static site.
- Backend logic needed for public data lookup/recommendation is moved to **Supabase Edge Functions**.
- Data/auth remain on Supabase Postgres + Supabase Auth.

## What changed for Supabase hosting compatibility
- Added `supabase/config.toml` with edge function config.
- Added edge functions:
  - `food-search`
  - `barcode-lookup`
  - `calorie-recommendation`
- Frontend now calls Supabase Edge Functions instead of requiring the separate Go API for these flows.

## Deploy from Git repo to Supabase
1. Push this repo to GitHub/GitLab.
2. In Supabase dashboard, connect the repository for web hosting.
3. Set frontend env vars:
   - `VITE_SUPABASE_URL`
   - `VITE_SUPABASE_ANON_KEY`
4. Deploy edge functions:
   ```bash
   supabase login
   supabase link --project-ref <your-project-ref>
   supabase functions deploy food-search
   supabase functions deploy barcode-lookup
   supabase functions deploy calorie-recommendation
   ```
5. Run migration in Supabase SQL editor: `supabase/migrations.sql`.

## Auth and verification
- Email/password signup with email confirmation.
- Phone OTP sign in / verification.
- Google OAuth.
All provided through Supabase Auth free tier (subject to provider limits).

## Notes
- `backend/` Go API remains in repo for optional self-hosted API extensions, but the core hosted path on Supabase no longer depends on it.
