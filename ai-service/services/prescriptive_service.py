import math
from typing import Dict, List, Optional
import numpy as np

class PrescriptiveService:
    """Heuristic prescriptive analytics for manufacturing dashboards."""

    def recommend_sales_actions(self, predictions: List[Dict], historical_avg: Optional[float] = None) -> List[Dict]:
        recommendations = []
        baseline = historical_avg or np.mean([p.get("recent_average", p.get("current_value", 0)) for p in predictions if p.get("current_value")]) or 0
        for item in predictions:
            forecast = item.get("next_month", 0)
            growth = item.get("growth_rate", 0)
            entity_name = item.get("entity_name") or f"Entity {item.get('entity_id')}"
            risk_level = "medium"
            action = "Hold course"
            rationale = "Stable outlook versus baseline performance."

            if growth < -5:
                risk_level = "high"
                action = "Deploy recovery plan"
                rationale = "Forecast indicates decline. Incentivise distributors and push targeted campaigns."
            elif growth < 0:
                risk_level = "medium"
                action = "Boost channel enablement"
                rationale = "Mild slowdown expected. Adjust pricing bundles and strengthen partner support."
            elif growth > 8:
                risk_level = "low"
                action = "Scale fulfilment"
                rationale = "Strong upside. Secure inventory and ensure logistics capacity."

            if baseline > 0:
                lift_vs_baseline = ((forecast - baseline) / baseline) * 100
            else:
                lift_vs_baseline = 0

            recommendations.append({
                "entity_id": item.get("entity_id"),
                "entity_name": entity_name,
                "predicted_sales": round(forecast, 2),
                "growth_rate": round(growth, 2),
                "lift_vs_baseline": round(lift_vs_baseline, 2),
                "risk_level": risk_level,
                "recommended_action": action,
                "rationale": rationale,
            })
        return recommendations

    def optimise_inventory(self, forecast_rows: List[Dict], current_stock: List[Dict]) -> Dict:
        lookup_stock = {row.get("gl_id"): row.get("on_hand_value", 0) for row in current_stock}
        recommendations = []
        slow_movers = []
        for row in forecast_rows:
            entity_id = row.get("entity_id") or row.get("product_id")
            demand = row.get("average_daily_forecast", 0) * row.get("forecast_period_days", 30) / 30
            reorder_point = demand * 0.35
            safety_stock = demand * 0.2
            on_hand = lookup_stock.get(entity_id, demand * 0.5)
            coverage_days = (on_hand / (row.get("average_daily_forecast", 1) + 1e-6))
            recommended_qty = max(0, reorder_point + safety_stock - on_hand)
            action = "Maintain" if recommended_qty <= 0 else "Raise purchase order"
            if coverage_days > 90:
                slow_movers.append({
                    "entity_id": entity_id,
                    "on_hand": round(on_hand, 2),
                    "coverage_days": round(coverage_days, 1),
                    "suggestion": "Investigate liquidation or transfer to high-demand channels"
                })
            recommendations.append({
                "entity_id": entity_id,
                "projected_demand": round(demand, 2),
                "on_hand": round(on_hand, 2),
                "reorder_point": round(reorder_point, 2),
                "safety_stock": round(safety_stock, 2),
                "coverage_days": round(coverage_days, 1),
                "recommended_order_qty": round(recommended_qty, 2),
                "action": action,
            })
        return {
            "inventory_actions": recommendations,
            "slow_movers": slow_movers,
        }

    def recommend_financial_actions(self, forecast_rows: List[Dict], expense_breakdown: List[Dict]) -> Dict:
        category_pressure = {}
        for expense in expense_breakdown:
            key = expense.get("expense_id")
            variance = (expense.get("actual_amount", 0) - expense.get("budget_amount", 0))
            category_pressure.setdefault(key, {"variance": 0, "periods": 0})
            category_pressure[key]["variance"] += variance
            category_pressure[key]["periods"] += 1
        insights = []
        for row in forecast_rows:
            entity_id = row.get("entity_id")
            avg_variance = 0
            if entity_id in category_pressure and category_pressure[entity_id]["periods"]:
                avg_variance = category_pressure[entity_id]["variance"] / category_pressure[entity_id]["periods"]
            confidence = row.get("confidence_level", 75)
            action = "Maintain spend discipline"
            rationale = "Spend is broadly aligned with targets."
            if avg_variance > 0:
                action = "Tighten discretionary spend"
                rationale = "Run-rate is above budget. Enforce approval gates and renegotiate suppliers."
            elif avg_variance < -50000:
                action = "Deploy surplus capital"
                rationale = "Under budget. Consider accelerating capex or invest surplus cash."
            insights.append({
                "entity_id": entity_id,
                "projected_total": round(row.get("total_forecast", 0), 2),
                "average_daily": round(row.get("average_daily_forecast", 0), 2),
                "confidence": confidence,
                "average_variance": round(avg_variance, 2),
                "recommended_action": action,
                "rationale": rationale,
            })
        return {
            "financial_actions": insights,
        }

    def enrich_anomalies(self, anomalies: Dict) -> Dict:
        enriched = []
        for anomaly in anomalies.get("anomalies", []):
            value = anomaly.get("value", 0)
            severity = "medium"
            recommendation = "Review related transactions"
            if value == 0:
                severity = "low"
                recommendation = "Validate data completeness"
            elif value > 0 and abs(value) > anomalies.get("mean", 0) * 1.5:
                severity = "high"
                recommendation = "Escalate to finance controller for approval check"
            enriched.append({
                **anomaly,
                "severity": severity,
                "recommended_action": recommendation,
            })
        anomalies["anomalies"] = enriched
        return anomalies

    def simulate_what_if(self, base_metrics: Dict, adjustments: Dict, horizon: int) -> Dict:
        sales = base_metrics.get("sales", 0)
        margin = base_metrics.get("gross_margin", 0)
        price_change = adjustments.get("price_change_pct", 0) / 100
        volume_change = adjustments.get("volume_change_pct", 0) / 100
        cost_change = adjustments.get("cost_change_pct", 0) / 100

        assumed_elasticity = -1.2
        adjusted_volume = volume_change if volume_change else price_change * assumed_elasticity
        projected_sales = sales * (1 + price_change + adjusted_volume)
        projected_margin = margin * (1 + price_change - cost_change)
        incremental_profit = (projected_sales * 0.1) - (sales * 0.1)

        return {
            "horizon": horizon,
            "inputs": adjustments,
            "projected_sales": round(projected_sales, 2),
            "projected_margin": round(projected_margin, 2),
            "incremental_profit": round(incremental_profit, 2),
            "narrative": "Scenario models price/volume elasticity to estimate revenue and margin shift.",
        }
