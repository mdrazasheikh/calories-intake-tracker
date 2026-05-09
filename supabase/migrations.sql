-- Supabase PostgreSQL schema

create table if not exists intake_entries (
  id bigint generated always as identity primary key,
  user_id uuid not null references auth.users(id) on delete cascade,
  food_name text not null,
  calories numeric not null check (calories >= 0),
  quantity numeric not null default 1 check (quantity > 0),
  serving_size text,
  consumed_at timestamptz not null default now()
);

create table if not exists daily_goals (
  user_id uuid primary key references auth.users(id) on delete cascade,
  target_calories integer not null check (target_calories > 0),
  protein_grams integer,
  carbs_grams integer,
  fat_grams integer,
  activity_level text not null,
  goal_description text,
  updated_at timestamptz not null default now()
);

create table if not exists weight_goals (
  user_id uuid primary key references auth.users(id) on delete cascade,
  current_weight_kg numeric not null check (current_weight_kg > 0),
  target_weight_kg numeric not null check (target_weight_kg > 0),
  weekly_change_kg numeric not null check (weekly_change_kg >= 0),
  goal_type text not null check (goal_type in ('lose','maintain','gain')),
  target_daily_kcal integer not null,
  estimated_days_goal integer not null,
  updated_at timestamptz not null default now()
);

create table if not exists user_security_settings (
  user_id uuid primary key references auth.users(id) on delete cascade,
  two_factor_enabled boolean not null default false,
  two_factor_method text,
  phone_verified boolean not null default false,
  email_verified boolean not null default false,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

alter table intake_entries enable row level security;
alter table daily_goals enable row level security;
alter table weight_goals enable row level security;
alter table user_security_settings enable row level security;

create policy "intake_entries_owner_rw" on intake_entries for all using (auth.uid() = user_id) with check (auth.uid() = user_id);
create policy "daily_goals_owner_rw" on daily_goals for all using (auth.uid() = user_id) with check (auth.uid() = user_id);
create policy "weight_goals_owner_rw" on weight_goals for all using (auth.uid() = user_id) with check (auth.uid() = user_id);
create policy "security_settings_owner_rw" on user_security_settings for all using (auth.uid() = user_id) with check (auth.uid() = user_id);
