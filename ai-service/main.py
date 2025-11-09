from fastapi import FastAPI, HTTPException, Query
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel
from typing import List, Optional, Dict, Any
import uvicorn
from datetime import datetime, timedelta
import os
from dotenv import load_dotenv

from services.forecast_service import ForecastService
from services.analysis_service import AnalysisService
from services.prescriptive_service import PrescriptiveService
from database.database import Database

load_dotenv()

app = FastAPI(
    title="AI Forecasting Service API",
    version="2.1.0",
    description="Advanced AI-powered forecasting and prescriptive analytics service for manufacturing industry"
)

# CORS middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Initialize services
forecast_service = ForecastService()
analysis_service = AnalysisService()
prescriptive_service = PrescriptiveService()
database = Database()

class ForecastRequest(BaseModel):
    forecast_type: str  # sales, production, finance, inventory
    period: str = "daily"  # daily, weekly, monthly
    start_date: str
    end_date: str
    days: int = 30
    company_id: Optional[int] = None
    factory_id: Optional[int] = None
    channel_id: Optional[int] = None
    product_id: Optional[int] = None
    budget_category: Optional[str] = None
    aggregation: Optional[str] = None

class PredictRequest(BaseModel):
    start_date: str
    end_date: str
    horizon: int = 30
    granularity: str = "channel"
    company_id: Optional[int] = None

class PrescriptiveRequest(BaseModel):
    module: str
    start_date: str
    end_date: str
    horizon: int = 30
    company_id: Optional[int] = None
    channel_id: Optional[int] = None
    product_id: Optional[int] = None
    budget_category: Optional[str] = None
    constraints: Optional[Dict[str, Any]] = None

class ScenarioRequest(BaseModel):
    horizon: int = 90
    base_metrics: Dict[str, float]
    adjustments: Dict[str, float]

class AnalysisRequest(BaseModel):
    metric: str
    start_date: str
    end_date: str
    analysis_type: str  # trend, anomaly, correlation

class SalesPipelineStagePayload(BaseModel):
    status: str
    orders: int = 0
    value: float = 0
    delivered_value: float = 0
    pending_value: float = 0
    avg_discount: float = 0
    avg_margin: float = 0
    avg_age_days: float = 0

class SalesPipelineSnapshotPayload(BaseModel):
    total_orders: int = 0
    total_value: float = 0
    delivered_value: float = 0
    pending_value: float = 0
    conversion_rate: float = 0
    stages: List[SalesPipelineStagePayload] = []

class SalesTargetPayload(BaseModel):
    channel_id: int
    channel_name: str
    revenue_target: float = 0
    actual_revenue: float = 0
    achievement: float = 0
    revenue_gap: float = 0
    promotion_budget: float = 0
    gross_margin_target: float = 0
    volume_target: float = 0
    new_customer_target: float = 0
    owner: Optional[str] = None

class SalesPromotionPayload(BaseModel):
    campaign_code: str
    campaign_name: str
    channel_name: Optional[str] = None
    spend_amount: float = 0
    revenue_uplift: float = 0
    uplift_percentage: float = 0
    roi: float = 0
    start_date: Optional[str] = None
    end_date: Optional[str] = None
    audience_tags: Optional[List[str]] = []

class SalesInsightRequest(BaseModel):
    targets: List[SalesTargetPayload] = []
    pipeline: SalesPipelineSnapshotPayload
    promotions: List[SalesPromotionPayload] = []


class ProductionInsightRequest(BaseModel):
    kpis: List[Dict[str, Any]] = []
    lines: List[Dict[str, Any]] = []
    wastage: List[Dict[str, Any]] = []
    maintenance: List[Dict[str, Any]] = []
    trend: List[Dict[str, Any]] = []
    forecast: Optional[Dict[str, Any]] = None


class FinanceInsightRequest(BaseModel):
    kpis: List[Dict[str, Any]] = []
    departments: List[Dict[str, Any]] = []
    categories: List[Dict[str, Any]] = []
    loans: List[Dict[str, Any]] = []
    trend: List[Dict[str, Any]] = []
    forecast: Optional[Dict[str, Any]] = None
    prescriptions: Optional[Dict[str, Any]] = None
    scenario: Optional[Dict[str, Any]] = None
    alerts: List[str] = []


class HRInsightRequest(BaseModel):
    kpis: List[Dict[str, Any]] = []
    departments: List[Dict[str, Any]] = []
    movements: List[Dict[str, Any]] = []
    trend: List[Dict[str, Any]] = []
    forecast: Optional[Dict[str, Any]] = None
    alerts: List[str] = []


class SupplyChainInsightRequest(BaseModel):
    kpis: List[Dict[str, Any]] = []
    suppliers: List[Dict[str, Any]] = []
    pending_orders: List[Dict[str, Any]] = []
    trend: List[Dict[str, Any]] = []
    forecast: Optional[Dict[str, Any]] = None
    alerts: List[str] = []


class InventoryInsightRequest(BaseModel):
    kpis: List[Dict[str, Any]] = []
    categories: List[Dict[str, Any]] = []
    companies: List[Dict[str, Any]] = []
    turnover: Optional[Dict[str, Any]] = None
    trend: List[Dict[str, Any]] = []
    forecast: Optional[Dict[str, Any]] = None
    prescriptions: Optional[Dict[str, Any]] = None
    slow_movers: Optional[List[Any]] = None
    alerts: List[str] = []

@app.get("/")
def read_root():
    return {
        "message": "AI Forecasting Service API",
        "version": "2.1.0",
        "endpoints": {
            "sales_forecast": "/api/forecast/sales",
            "production_forecast": "/api/forecast/production",
            "financial_forecast": "/api/forecast/finance",
            "inventory_forecast": "/api/forecast/inventory",
            "sales_prediction": "/api/predict/sales/summary",
            "inventory_prescription": "/api/prescribe/inventory",
            "financial_prescription": "/api/prescribe/finance",
            "scenario": "/api/scenario/whatif",
            "analysis": "/api/analyze"
        }
    }

@app.post("/api/forecast/sales")
async def create_sales_forecast(request: ForecastRequest):
    """Generate sales forecast with AI"""
    try:
        # Get historical data from database
        data = database.get_sales_data(
            start_date=request.start_date,
            end_date=request.end_date,
            company_id=request.company_id,
            channel_id=request.channel_id,
            product_id=request.product_id
        )
        
        if not data or len(data) == 0:
            raise HTTPException(status_code=404, detail="No sales data found for the specified period")
        
        # Generate forecast
        forecast = forecast_service.generate_sales_forecast(
            data=data,
            days=request.days,
            company_id=request.company_id,
            channel_id=request.channel_id,
            product_id=request.product_id
        )
        
        # Store forecast in database
        forecast_id = database.save_forecast(forecast)
        
        forecast["forecast_id"] = forecast_id
        return forecast
        
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.post("/api/forecast/production")
async def create_production_forecast(request: ForecastRequest):
    """Generate production forecast with AI"""
    try:
        data = database.get_production_data(
            start_date=request.start_date,
            end_date=request.end_date,
            factory_id=request.factory_id,
            product_id=request.product_id
        )
        
        if not data or len(data) == 0:
            raise HTTPException(status_code=404, detail="No production data found")
        
        forecast = forecast_service.generate_production_forecast(
            data=data,
            days=request.days,
            factory_id=request.factory_id,
            product_id=request.product_id
        )
        
        forecast_id = database.save_forecast(forecast)
        forecast["forecast_id"] = forecast_id
        return forecast
        
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.post("/api/forecast/finance")
async def create_financial_forecast(request: ForecastRequest):
    """Generate financial forecast with AI"""
    try:
        data = database.get_financial_data(
            start_date=request.start_date,
            end_date=request.end_date,
            budget_category=request.budget_category
        )
        
        if not data or len(data) == 0:
            raise HTTPException(status_code=404, detail="No financial data found")
        
        forecast = forecast_service.generate_financial_forecast(
            data=data,
            days=request.days,
            budget_category=request.budget_category
        )
        
        forecast_id = database.save_forecast(forecast)
        forecast["forecast_id"] = forecast_id
        return forecast
        
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.post("/api/forecast/inventory")
async def create_inventory_forecast(request: ForecastRequest):
    """Generate inventory forecast with AI"""
    try:
        data = database.get_inventory_data(
            start_date=request.start_date,
            end_date=request.end_date,
            product_id=request.product_id
        )
        
        if not data or len(data) == 0:
            raise HTTPException(status_code=404, detail="No inventory data found")
        
        forecast = forecast_service.generate_inventory_forecast(
            data=data,
            days=request.days,
            product_id=request.product_id
        )
        
        forecast_id = database.save_forecast(forecast)
        forecast["forecast_id"] = forecast_id
        return forecast
        
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.post("/api/analyze")
async def analyze_data(request: AnalysisRequest):
    """Perform data analysis (trend, anomaly, correlation)"""
    try:
        data = database.get_historical_data(
            metric=request.metric,
            start_date=request.start_date,
            end_date=request.end_date
        )
        
        if not data or len(data) == 0:
            raise HTTPException(status_code=404, detail="No data found for analysis")
        
        if request.analysis_type == "trend":
            result = analysis_service.analyze_trend(data)
        elif request.analysis_type == "anomaly":
            result = analysis_service.detect_anomalies(data)
        elif request.analysis_type == "correlation":
            result = analysis_service.analyze_correlation(data)
        else:
            raise HTTPException(status_code=400, detail="Invalid analysis type")
        
        return {
            "analysis_type": request.analysis_type,
            "metric": request.metric,
            "result": result,
            "created_at": datetime.now().isoformat()
        }
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.post("/api/predict/sales/summary")
async def predict_sales_summary(request: PredictRequest):
    try:
        data = database.get_sales_breakdown(request.start_date, request.end_date, request.granularity)
        if not data:
            raise HTTPException(status_code=404, detail="No historical sales data available")
        predictions = forecast_service.predict_grouped_sales(data, request.granularity, request.horizon)
        recommendations = prescriptive_service.recommend_sales_actions(predictions)
        return {
            "granularity": request.granularity,
            "horizon": request.horizon,
            "predictions": predictions,
            "recommendations": recommendations,
            "generated_at": datetime.now().isoformat()
        }
    except HTTPException:
        raise
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


@app.post("/api/enrich/sales/insights")
async def enrich_sales_insights(request: SalesInsightRequest):
    try:
        pipeline_summary = prescriptive_service.summarise_pipeline(request.pipeline.dict())
        targets_summary = prescriptive_service.assess_targets([target.dict() for target in request.targets])
        promotions_summary = prescriptive_service.analyse_promotions([promo.dict() for promo in request.promotions])
        executive_summary = prescriptive_service.compile_sales_executive_summary(pipeline_summary, targets_summary, promotions_summary)
        recommended_actions = prescriptive_service.propose_sales_actions(pipeline_summary, targets_summary, promotions_summary)
        return {
            "pipeline": pipeline_summary,
            "targets": targets_summary,
            "promotions": promotions_summary,
            "executive_summary": executive_summary,
            "recommended_actions": recommended_actions,
            "generated_at": datetime.now().isoformat()
        }
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))
 
 
@app.post("/api/enrich/production/insights")
async def enrich_production_insights(request: ProductionInsightRequest):
    try:
        summary = prescriptive_service.summarise_production_overview(request.dict())
        summary["generated_at"] = datetime.now().isoformat()
        return summary
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


@app.post("/api/enrich/finance/insights")
async def enrich_finance_insights(request: FinanceInsightRequest):
    try:
        summary = prescriptive_service.summarise_finance_overview(request.dict())
        summary["generated_at"] = datetime.now().isoformat()
        return summary
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))
 
 
@app.post("/api/enrich/hr/insights")
async def enrich_hr_insights(request: HRInsightRequest):
    try:
        summary = prescriptive_service.summarise_hr_overview(request.dict())
        summary["generated_at"] = datetime.now().isoformat()
        return summary
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


@app.post("/api/enrich/supplychain/insights")
async def enrich_supplychain_insights(request: SupplyChainInsightRequest):
    try:
        summary = prescriptive_service.summarise_supply_chain_overview(request.dict())
        summary["generated_at"] = datetime.now().isoformat()
        return summary
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


@app.post("/api/enrich/inventory/insights")
async def enrich_inventory_insights(request: InventoryInsightRequest):
    try:
        summary = prescriptive_service.summarise_inventory_overview(request.dict())
        summary["generated_at"] = datetime.now().isoformat()
        return summary
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


@app.post("/api/prescribe/inventory")
async def prescribe_inventory(request: PrescriptiveRequest):
    try:
        history = database.get_inventory_data(request.start_date, request.end_date, request.product_id)
        if not history:
            raise HTTPException(status_code=404, detail="No inventory history available")
        forecast = forecast_service.generate_inventory_forecast(history, days=request.horizon, product_id=request.product_id)
        snapshot = database.get_current_inventory_snapshot()
        prescriptions = prescriptive_service.optimise_inventory([forecast], snapshot)
        return {
            "forecast": forecast,
            "prescriptions": prescriptions,
            "generated_at": datetime.now().isoformat()
        }
    except HTTPException:
        raise
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.post("/api/prescribe/finance")
async def prescribe_finance(request: PrescriptiveRequest):
    try:
        financial_history = database.get_financial_data(request.start_date, request.end_date, request.budget_category)
        if not financial_history:
            raise HTTPException(status_code=404, detail="No financial data available")
        forecast = forecast_service.generate_financial_forecast(financial_history, days=request.horizon, budget_category=request.budget_category)
        expense_breakdown = database.get_expense_breakdown(request.start_date, request.end_date)
        prescriptions = prescriptive_service.recommend_financial_actions([forecast], expense_breakdown)
        return {
            "forecast": forecast,
            "prescriptions": prescriptions,
            "generated_at": datetime.now().isoformat()
        }
    except HTTPException:
        raise
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.post("/api/analyze/anomaly/enriched")
async def analyze_anomalies(request: AnalysisRequest):
    try:
        data = database.get_historical_data(request.metric, request.start_date, request.end_date)
        if not data:
            raise HTTPException(status_code=404, detail="No data found for analysis")
        anomalies = analysis_service.detect_anomalies(data)
        enriched = prescriptive_service.enrich_anomalies(anomalies)
        enriched.update({
            "metric": request.metric,
            "analysis_type": "anomaly",
            "generated_at": datetime.now().isoformat()
        })
        return enriched
    except HTTPException:
        raise
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.post("/api/scenario/whatif")
async def simulate_scenario(request: ScenarioRequest):
    try:
        result = prescriptive_service.simulate_what_if(request.base_metrics, request.adjustments, request.horizon)
        return result
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.get("/api/forecast/{forecast_id}")
async def get_forecast(forecast_id: str):
    """Get forecast by ID"""
    try:
        forecast = database.get_forecast(forecast_id)
        if not forecast:
            raise HTTPException(status_code=404, detail="Forecast not found")
        return forecast
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.get("/health")
def health_check():
    return {"status": "healthy", "service": "AI Forecasting Service"}

if __name__ == "__main__":
    port = int(os.getenv("PORT", 8000))
    uvicorn.run(app, host="0.0.0.0", port=port)
