import { serve } from 'https://deno.land/std@0.224.0/http/server.ts'

serve(async (req) => {
  const body = await req.json()
  const { sex, age, height_cm, weight_kg, activity_level, goal_type } = body
  let bmr = 10 * weight_kg + 6.25 * height_cm - 5 * age
  bmr += sex === 'male' ? 5 : -161
  const multipliers: Record<string, number> = { sedentary: 1.2, light: 1.375, moderate: 1.55, active: 1.725, very_active: 1.9 }
  const maintenance = bmr * (multipliers[activity_level] ?? 1.2)
  let recommended = maintenance
  let guidance = 'Maintain current weight by staying near maintenance calories.'
  if (goal_type === 'lose') { recommended = Math.max(1200, maintenance - 500); guidance = 'Target ~500 kcal/day deficit.' }
  if (goal_type === 'gain') { recommended = maintenance + 300; guidance = 'Target ~300 kcal/day surplus.' }
  return Response.json({ bmr, maintenance_calories: maintenance, recommended_calories: recommended, goal_type, guidance })
})
