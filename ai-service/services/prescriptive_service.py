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

    def summarise_pipeline(self, pipeline: Dict) -> Dict:
        stages = pipeline.get("stages", []) or []
        total_value = float(pipeline.get("total_value", 0) or 0)
        delivered_value = float(pipeline.get("delivered_value", 0) or 0)
        conversion_rate = pipeline.get("conversion_rate")
        if conversion_rate is None:
            conversion_rate = (delivered_value / total_value * 100) if total_value else 0

        alerts = []
        recommendations = []
        for stage in stages:
            status = (stage.get("status") or "").lower()
            orders = int(stage.get("orders") or 0)
            pending_value = float(stage.get("pending_value", stage.get("value", 0)) or 0)
            avg_age = float(stage.get("avg_age_days") or 0)
            avg_discount = float(stage.get("avg_discount") or 0)

            if status in {"draft", "confirmed"} and avg_age > 9 and pending_value > 0:
                alerts.append(f"{orders} {status} orders ageing {avg_age:.1f} days (৳{pending_value:,.0f})")
                recommendations.append({
                    "category": "pipeline",
                    "status": status,
                    "orders": orders,
                    "value": round(pending_value, 2),
                    "action": "Escalate aged orders for decision",
                    "impact": "Release stuck revenue"
                })
            if status == "dispatching" and avg_age > 5 and pending_value > 0:
                recommendations.append({
                    "category": "pipeline",
                    "status": status,
                    "orders": orders,
                    "value": round(pending_value, 2),
                    "action": "Prioritise logistics slot and shipping",
                    "impact": "Prevent slippage in fulfilment"
                })
            if avg_discount > 8 and pending_value > 0:
                alerts.append(f"{status.title()} stage discount averaging {avg_discount:.1f}%")

        narrative = f"Pipeline conversion {conversion_rate:.1f}% on ৳{total_value:,.0f} order book" if total_value else "Pipeline awaiting new bookings"

        return {
            "narrative": narrative,
            "alerts": alerts,
            "stage_recommendations": recommendations,
            "conversion_rate": round(conversion_rate, 2),
            "total_value": round(total_value, 2),
        }

    def assess_targets(self, targets: List[Dict]) -> Dict:
        alerts = []
        focus = []
        for target in targets:
            channel = target.get("channel_name") or "Channel"
            gap = float(target.get("revenue_gap") or 0)
            target_value = float(target.get("revenue_target") or 0)
            achievement = float(target.get("achievement") or 0)
            if target_value <= 0:
                continue
            gap_pct = (gap / target_value * 100) if target_value else 0
            if gap > 0 and gap_pct >= 12:
                alerts.append(f"{channel} ({achievement:.1f}%) lags target by ৳{gap:,.0f}")
                focus.append({
                    "category": "target",
                    "channel": channel,
                    "gap": round(gap, 2),
                    "gap_pct": round(gap_pct, 2),
                    "action": "Launch bundle promo + distributor incentive",
                    "impact": "Recover revenue gap",
                })
        return {"alerts": alerts, "focus": focus}

    def analyse_promotions(self, promotions: List[Dict]) -> Dict:
        highlights = []
        underperformers = []
        for promo in promotions:
            name = promo.get("campaign_name") or promo.get("campaign_code") or "Campaign"
            roi = float(promo.get("roi") or 0)
            uplift = float(promo.get("revenue_uplift") or 0)
            spend = float(promo.get("spend_amount") or 0)
            if roi >= 150 and uplift > 0:
                highlights.append(f"{name} ROI {roi:.0f}% delivering ৳{uplift:,.0f} uplift")
            if spend > 0 and roi < 90:
                underperformers.append({
                    "category": "promotion",
                    "campaign": name,
                    "roi": round(roi, 2),
                    "action": "Tweak targeting / pause low ROI slots",
                    "impact": "Reallocate spend to better performers",
                })
        return {"highlights": highlights, "underperformers": underperformers}

    def compile_sales_executive_summary(self, pipeline: Dict, targets: Dict, promotions: Dict) -> List[str]:
        summary = []
        narrative = pipeline.get("narrative")
        if narrative:
            summary.append(narrative)
        if targets.get("alerts"):
            summary.append(f"{len(targets['alerts'])} channel gap alerts raised")
        if promotions.get("highlights"):
            summary.append(promotions["highlights"][0])
        return summary

    def propose_sales_actions(self, pipeline: Dict, targets: Dict, promotions: Dict) -> List[Dict]:
        actions: List[Dict] = []
        actions.extend(pipeline.get("stage_recommendations", [])[:3])
        actions.extend(targets.get("focus", [])[:2])
        actions.extend(promotions.get("underperformers", [])[:2])
        return actions

    def summarise_production_overview(self, payload: Dict) -> Dict:
        lines = payload.get("lines") or []
        wastage = payload.get("wastage") or []
        maintenance = payload.get("maintenance") or []

        executive_summary: List[str] = []
        insights: List[str] = []
        recommendations: List[Dict] = []

        if lines:
            sorted_lines = sorted(lines, key=lambda item: item.get("oee", item.get("efficiency", 0)), reverse=True)
            top_line = sorted_lines[0]
            executive_summary.append(
                f"Top line OEE {top_line.get('oee', top_line.get('efficiency', 0)):.1f}% with output {top_line.get('actual_output', 0):,.0f}"
            )

            lagging = sorted(lines, key=lambda item: item.get("efficiency", 0))[:1]
            for entry in lagging:
                efficiency = entry.get('efficiency', 0)
                downtime = entry.get('downtime_minutes', 0)
                insights.append(
                    f"Line {entry.get('production_line_id', entry.get('factory_id'))} efficiency {efficiency:.1f}% with {downtime:,.0f} mins downtime"
                )
                if efficiency < 75 or downtime > 600:
                    recommendations.append({
                        "category": "line",
                        "line_id": entry.get('production_line_id'),
                        "action": "Run SMED + bottleneck audit",
                        "impact": "Recover lost throughput",
                    })

        if wastage:
            worst = max(wastage, key=lambda item: item.get('rate', 0))
            insights.append(
                f"Highest wastage at factory {worst.get('factory')} with rate {worst.get('rate', 0):.1f}% (৳ {worst.get('amount', 0):,.0f})"
            )
            if worst.get('rate', 0) > 5:
                recommendations.append({
                    "category": "wastage",
                    "factory": worst.get('factory'),
                    "action": "Deploy quality circle on scrap hotspots",
                    "impact": "Reduce scrap cost",
                })

        if maintenance:
            maintenance_sorted = sorted(maintenance, key=lambda item: item.get('downtime', 0), reverse=True)
            top_issue = maintenance_sorted[0]
            insights.append(
                f"Machine {top_issue.get('machine_code')} caused {top_issue.get('downtime', 0):,.0f} mins downtime across {top_issue.get('events', 0)} events"
            )
            if top_issue.get('events', 0) > 2:
                recommendations.append({
                    "category": "maintenance",
                    "machine": top_issue.get('machine_code'),
                    "action": "Schedule root-cause teardown and preventive maintenance window",
                    "impact": "Stabilise uptime",
                })

        if not executive_summary and lines:
            total_output = sum(item.get('actual_output', 0) for item in lines)
            executive_summary.append(f"Total output {total_output:,.0f} units produced this period")

        return {
            "executive_summary": executive_summary,
            "insights": insights,
            "recommended_actions": recommendations,
        }

    def summarise_finance_overview(self, payload: Dict) -> Dict:
        kpis = payload.get("kpis") or []
        departments = payload.get("departments") or []
        categories = payload.get("categories") or []
        loans = payload.get("loans") or []

        executive_summary: List[str] = []
        insights: List[str] = []
        recommendations: List[Dict] = []

        variance = next((kpi for kpi in kpis if kpi.get('label') == 'Variance'), None)
        if variance:
            value = variance.get('value', 0)
            executive_summary.append(f"Budget variance currently ৳ {value:,.0f}")

        if departments:
            worst = min(departments, key=lambda item: item.get('variance_percent', 0))
            insights.append(
                f"Dept {worst.get('department_name')} variance {worst.get('variance_percent', 0):.1f}% ({worst.get('variance', 0):,.0f})"
            )
            if worst.get('variance_percent', 0) < -10:
                recommendations.append({
                    "category": "department",
                    "department": worst.get('department_name'),
                    "action": "Freeze discretionary spend & review vendor contracts",
                    "impact": "Bring variance within 5%",
                })

        if categories:
            top_category = max(categories, key=lambda item: item.get('actual', 0))
            insights.append(
                f"Highest spend category: {top_category.get('category_name')} at ৳ {top_category.get('actual', 0):,.0f}"
            )

        if loans:
            largest = max(loans, key=lambda item: item.get('amount', 0))
            insights.append(
                f"Largest loan head {largest.get('head')} totals ৳ {largest.get('amount', 0):,.0f}"
            )
            recommendations.append({
                "category": "liquidity",
                "loan_head": largest.get('head'),
                "action": "Evaluate refinancing to reduce financial expense",
                "impact": "Improve cash cost of capital",
            })

        return {
            "executive_summary": executive_summary,
            "insights": insights,
            "recommended_actions": recommendations,
        }