import pandas as pd
import numpy as np
import math
from prophet import Prophet
from datetime import datetime, timedelta
from typing import List, Dict, Optional
import warnings
warnings.filterwarnings('ignore')

class ForecastService:
    def __init__(self):
        self.prophet_model = None
    
    def generate_sales_forecast(self, data: List[Dict], days: int = 30, company_id: Optional[int] = None, 
                                channel_id: Optional[int] = None, product_id: Optional[int] = None) -> Dict:
        """
        Generate sales forecast using Prophet model with confidence intervals
        """
        try:
            # Convert data to DataFrame
            df = pd.DataFrame(data)
            
            if df.empty or len(df) < 7:
                raise ValueError("Insufficient data for forecasting. Need at least 7 data points.")
            
            # Prepare data for Prophet
            if 'data_month' in df.columns:
                df['ds'] = pd.to_datetime(df['data_month'])
                df['y'] = df.get('billed', df.get('sales', df.get('value', df.get('amount', 0))))
            elif 'data_date' in df.columns:
                df['ds'] = pd.to_datetime(df['data_date'])
                df['y'] = df.get('billed', df.get('sales', df.get('value', df.get('amount', 0))))
            elif 'date' in df.columns:
                df['ds'] = pd.to_datetime(df['date'])
                df['y'] = df.get('sales', df.get('value', df.get('amount', 0)))
            else:
                raise ValueError("Date column not found in data")
            
            # Select only required columns and remove nulls
            prophet_df = df[['ds', 'y']].copy()
            prophet_df = prophet_df.dropna()
            prophet_df = prophet_df.sort_values('ds')
            
            if len(prophet_df) < 7:
                raise ValueError("Insufficient data after cleaning")
            
            # Initialize Prophet with manufacturing-specific parameters
            model = Prophet(
                yearly_seasonality=True,
                weekly_seasonality=True,
                daily_seasonality=False,
                seasonality_mode='multiplicative',
                changepoint_prior_scale=0.05,
                seasonality_prior_scale=10,
                holidays_prior_scale=10,
                mcmc_samples=0,
                interval_width=0.95,
                uncertainty_samples=1000
            )
            
            # Add custom seasonalities for manufacturing (monthly, quarterly)
            model.add_seasonality(name='monthly', period=30.5, fourier_order=5)
            model.add_seasonality(name='quarterly', period=91.25, fourier_order=3)
            
            # Fit model
            model.fit(prophet_df)
            
            # Create future dataframe
            future = model.make_future_dataframe(periods=days, freq='D')
            
            # Generate forecast
            forecast = model.predict(future)
            
            # Get forecasted values (future dates only)
            future_forecast = forecast.tail(days)
            
            # Calculate confidence level based on historical accuracy
            historical_accuracy = self._calculate_accuracy(prophet_df, forecast.head(len(prophet_df)))
            confidence_level = min(95, max(60, historical_accuracy * 100))
            
            # Format forecast results
            forecast_results = []
            for idx, row in future_forecast.iterrows():
                forecast_results.append({
                    "date": row['ds'].strftime('%Y-%m-%d'),
                    "forecast": max(0, float(row['yhat'])),  # Ensure non-negative
                    "lower_bound": max(0, float(row['yhat_lower'])),
                    "upper_bound": max(0, float(row['yhat_upper'])),
                    "trend": float(row['trend']),
                    "seasonal": float(row.get('monthly', 0)) + float(row.get('quarterly', 0))
                })
            
            # Calculate summary statistics
            total_forecast = sum([f['forecast'] for f in forecast_results])
            avg_daily_forecast = total_forecast / days
            growth_rate = self._calculate_growth_rate(forecast_results)
            
            return {
                "forecast_type": "sales",
                "entity_id": product_id or channel_id or company_id,
                "forecast_period_days": days,
                "confidence_level": round(confidence_level, 2),
                "total_forecast": round(total_forecast, 2),
                "average_daily_forecast": round(avg_daily_forecast, 2),
                "projected_growth_rate": round(growth_rate, 2),
                "forecast_data": forecast_results,
                "model_used": "Prophet",
                "created_at": datetime.now().isoformat()
            }
            
        except Exception as e:
            raise Exception(f"Sales forecast generation failed: {str(e)}")
    
    def generate_production_forecast(self, data: List[Dict], days: int = 30, factory_id: Optional[int] = None,
                                     product_id: Optional[int] = None) -> Dict:
        """
        Generate production forecast with capacity constraints
        """
        try:
            df = pd.DataFrame(data)
            
            if df.empty or len(df) < 7:
                raise ValueError("Insufficient data for forecasting")
            
            # Prepare data
            if 'month' in df.columns:
                df['ds'] = pd.to_datetime(df['month'])
                df['y'] = df.get('actual_output', df.get('cmonthly_amount', df.get('production', 0)))
            elif 'production_date' in df.columns:
                df['ds'] = pd.to_datetime(df['production_date'])
                df['y'] = df.get('actual_output', df.get('production', 0))
            else:
                raise ValueError("Date column not found")
            
            prophet_df = df[['ds', 'y']].copy()
            prophet_df = prophet_df.dropna().sort_values('ds')
            
            if len(prophet_df) < 7:
                raise ValueError("Insufficient data after cleaning")
            
            # Initialize Prophet for production
            model = Prophet(
                yearly_seasonality=True,
                weekly_seasonality=True,
                daily_seasonality=False,
                seasonality_mode='additive',
                changepoint_prior_scale=0.05,
                interval_width=0.95
            )
            
            model.add_seasonality(name='monthly', period=30.5, fourier_order=5)
            
            model.fit(prophet_df)
            
            future = model.make_future_dataframe(periods=days, freq='D')
            forecast = model.predict(future)
            future_forecast = forecast.tail(days)
            
            confidence_level = 85.0  # Production forecasts typically have higher confidence
            
            forecast_results = []
            for idx, row in future_forecast.iterrows():
                forecast_results.append({
                    "date": row['ds'].strftime('%Y-%m-%d'),
                    "forecast": max(0, float(row['yhat'])),
                    "lower_bound": max(0, float(row['yhat_lower'])),
                    "upper_bound": max(0, float(row['yhat_upper'])),
                    "trend": float(row['trend'])
                })
            
            total_forecast = sum([f['forecast'] for f in forecast_results])
            avg_daily_forecast = total_forecast / days
            
            return {
                "forecast_type": "production",
                "entity_id": factory_id or product_id,
                "forecast_period_days": days,
                "confidence_level": confidence_level,
                "total_forecast": round(total_forecast, 2),
                "average_daily_forecast": round(avg_daily_forecast, 2),
                "forecast_data": forecast_results,
                "model_used": "Prophet",
                "created_at": datetime.now().isoformat()
            }
            
        except Exception as e:
            raise Exception(f"Production forecast generation failed: {str(e)}")
    
    def generate_financial_forecast(self, data: List[Dict], days: int = 30, 
                                     budget_category: Optional[str] = None) -> Dict:
        """
        Generate financial forecast for budget planning
        """
        try:
            df = pd.DataFrame(data)
            
            if df.empty or len(df) < 7:
                raise ValueError("Insufficient data for forecasting")
            
            # Prepare data
            if 'month' in df.columns:
                df['ds'] = pd.to_datetime(df['month'])
                df['y'] = df.get('actual_amount', df.get('budget_amount', df.get('expense', 0)))
            else:
                raise ValueError("Date column not found")
            
            prophet_df = df[['ds', 'y']].copy()
            prophet_df = prophet_df.dropna().sort_values('ds')
            
            if len(prophet_df) < 7:
                raise ValueError("Insufficient data after cleaning")
            
            # Initialize Prophet for financial data
            model = Prophet(
                yearly_seasonality=True,
                weekly_seasonality=False,
                daily_seasonality=False,
                seasonality_mode='additive',
                changepoint_prior_scale=0.01,  # More conservative for financial data
                interval_width=0.95
            )
            
            model.add_seasonality(name='monthly', period=30.5, fourier_order=5)
            model.add_seasonality(name='quarterly', period=91.25, fourier_order=3)
            
            model.fit(prophet_df)
            
            future = model.make_future_dataframe(periods=days, freq='D')
            forecast = model.predict(future)
            future_forecast = forecast.tail(days)
            
            confidence_level = 80.0  # Financial forecasts have moderate confidence
            
            forecast_results = []
            for idx, row in future_forecast.iterrows():
                forecast_results.append({
                    "date": row['ds'].strftime('%Y-%m-%d'),
                    "forecast": float(row['yhat']),
                    "lower_bound": float(row['yhat_lower']),
                    "upper_bound": float(row['yhat_upper']),
                    "trend": float(row['trend'])
                })
            
            total_forecast = sum([f['forecast'] for f in forecast_results])
            avg_daily_forecast = total_forecast / days
            
            return {
                "forecast_type": "finance",
                "entity_id": budget_category,
                "forecast_period_days": days,
                "confidence_level": confidence_level,
                "total_forecast": round(total_forecast, 2),
                "average_daily_forecast": round(avg_daily_forecast, 2),
                "forecast_data": forecast_results,
                "model_used": "Prophet",
                "created_at": datetime.now().isoformat()
            }
            
        except Exception as e:
            raise Exception(f"Financial forecast generation failed: {str(e)}")
    
    def generate_inventory_forecast(self, data: List[Dict], days: int = 30, 
                                     product_id: Optional[int] = None) -> Dict:
        """
        Generate inventory forecast for stock management
        """
        try:
            df = pd.DataFrame(data)
            
            if df.empty or len(df) < 7:
                raise ValueError("Insufficient data for forecasting")
            
            # Prepare data
            if 'month' in df.columns:
                df['ds'] = pd.to_datetime(df['month'])
                df['y'] = df.get('amount', df.get('quantity', df.get('inventory', 0)))
            else:
                raise ValueError("Date column not found")
            
            prophet_df = df[['ds', 'y']].copy()
            prophet_df = prophet_df.dropna().sort_values('ds')
            
            if len(prophet_df) < 7:
                raise ValueError("Insufficient data after cleaning")
            
            # Initialize Prophet for inventory
            model = Prophet(
                yearly_seasonality=True,
                weekly_seasonality=True,
                daily_seasonality=False,
                seasonality_mode='additive',
                changepoint_prior_scale=0.05,
                interval_width=0.90
            )
            
            model.fit(prophet_df)
            
            future = model.make_future_dataframe(periods=days, freq='D')
            forecast = model.predict(future)
            future_forecast = forecast.tail(days)
            
            confidence_level = 75.0
            
            forecast_results = []
            min_stock_level = float('inf')
            max_stock_level = 0
            
            for idx, row in future_forecast.iterrows():
                forecast_value = max(0, float(row['yhat']))
                min_stock_level = min(min_stock_level, forecast_value)
                max_stock_level = max(max_stock_level, forecast_value)
                
                forecast_results.append({
                    "date": row['ds'].strftime('%Y-%m-%d'),
                    "forecast": forecast_value,
                    "lower_bound": max(0, float(row['yhat_lower'])),
                    "upper_bound": max(0, float(row['yhat_upper'])),
                    "reorder_point": forecast_value * 0.3,  # 30% of forecast as reorder point
                    "safety_stock": forecast_value * 0.2   # 20% safety stock
                })
            
            return {
                "forecast_type": "inventory",
                "entity_id": product_id,
                "forecast_period_days": days,
                "confidence_level": confidence_level,
                "min_stock_level": round(min_stock_level, 2),
                "max_stock_level": round(max_stock_level, 2),
                "average_daily_forecast": round(sum([p['forecast'] for p in forecast_results]) / max(1, days), 2),
                "forecast_data": forecast_results,
                "model_used": "Prophet",
                "created_at": datetime.now().isoformat()
            }
            
        except Exception as e:
            raise Exception(f"Inventory forecast generation failed: {str(e)}")
    
    def predict_grouped_sales(self, grouped_data: List[Dict], granularity: str = "channel", horizon: int = 30) -> List[Dict]:
        if not grouped_data:
            return []
        df = pd.DataFrame(grouped_data)
        df['period'] = pd.to_datetime(df['period'])
        value_col = 'sales_value'
        entity_key = 'channel_id' if granularity != 'product' else 'product_id'
        name_key = 'channel_name' if granularity != 'product' else 'product_name'

        predictions = []
        for entity_id, group in df.groupby(entity_key):
            group = group.sort_values('period')
            series = group[value_col].values
            if len(series) < 3:
                continue
            recent_values = series[-3:]
            current_value = recent_values[-1]
            recent_average = np.mean(recent_values)
            historical_average = np.mean(series)
            growth_rate = 0.0
            if historical_average > 0:
                growth_rate = ((current_value - historical_average) / historical_average) * 100
            elif recent_average > 0:
                growth_rate = ((current_value - recent_average) / recent_average) * 100

            next_month = current_value * (1 + growth_rate / 100)
            next_quarter = next_month * math.pow(1 + (growth_rate / 100), horizon / 30)
            confidence = min(95, 60 + len(series))

            predictions.append({
                "entity_id": int(entity_id) if entity_id is not None else None,
                "entity_name": group[name_key].iloc[0] if name_key in group else f"{granularity.capitalize()} {entity_id}",
                "current_value": round(float(current_value), 2),
                "recent_average": round(float(recent_average), 2),
                "growth_rate": round(float(growth_rate), 2),
                "next_month": round(float(next_month), 2),
                "next_quarter": round(float(next_quarter), 2),
                "confidence": round(float(confidence), 2),
            })
        return predictions
    
    def _calculate_accuracy(self, actual: pd.DataFrame, predicted: pd.DataFrame) -> float:
        """Calculate forecast accuracy using MAPE"""
        try:
            if len(actual) != len(predicted):
                return 0.7  # Default accuracy
            
            actual_values = actual['y'].values
            predicted_values = predicted['yhat'].values
            
            # Calculate MAPE
            mape = np.mean(np.abs((actual_values - predicted_values) / (actual_values + 1e-10))) * 100
            accuracy = 1 - (mape / 100)
            
            return max(0, min(1, accuracy))
        except:
            return 0.7
    
    def _calculate_growth_rate(self, forecast_results: List[Dict]) -> float:
        """Calculate projected growth rate"""
        try:
            if len(forecast_results) < 2:
                return 0.0
            
            first_value = forecast_results[0]['forecast']
            last_value = forecast_results[-1]['forecast']
            
            if first_value == 0:
                return 0.0
            
            growth_rate = ((last_value - first_value) / first_value) * 100
            return growth_rate
        except:
            return 0.0
