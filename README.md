# Calories Intake Tracker (Supabase-only Architecture)

You're right: **Supabase does not natively host Go services**. This project is now fully aligned to Supabase-native runtime.

## Final Stack (Supabase native)
- **Frontend:** Vue 3 static site (deploy from Git repo).
- **Backend logic:** Supabase Edge Functions (Deno/TypeScript).
- **Database:** Supabase Postgres + RLS.
- **Auth:** Supabase Auth (email, phone OTP, Google OAuth).

## Supabase-hosted components in this repo
- `supabase/functions/food-search`
- `supabase/functions/barcode-lookup`
- `supabase/functions/calorie-recommendation`
- `supabase/migrations.sql`
- `supabase/config.toml`
- `frontend/` Vue app

## Deploy from Git repo
1. Push this repository to GitHub/GitLab.
2. In Supabase dashboard, connect your Git repo for the frontend hosting flow.
3. Configure frontend env vars:
   - `VITE_SUPABASE_URL`
   - `VITE_SUPABASE_ANON_KEY`
4. Deploy functions with Supabase CLI:
   ```bash
   supabase login
   supabase link --project-ref <project-ref>
   supabase functions deploy food-search
   supabase functions deploy barcode-lookup
   supabase functions deploy calorie-recommendation
   ```
5. Run SQL migration in Supabase SQL Editor:
   - `supabase/migrations.sql`

## Auth/verification setup
Enable in Supabase dashboard:
- Email provider with confirm email
- Phone provider with OTP
- Google OAuth provider

## Note
The previous optional Go backend has been removed so the codebase is now strictly compatible with Supabase-native hosting.
