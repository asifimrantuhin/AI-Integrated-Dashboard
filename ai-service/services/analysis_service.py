import pandas as pd
import numpy as np
from sklearn.ensemble import IsolationForest
from sklearn.preprocessing import StandardScaler
from typing import List, Dict

class AnalysisService:
    def __init__(self):
        self.isolation_forest = IsolationForest(contamination=0.1)
        self.scaler = StandardScaler()
    
    def analyze_trend(self, data: List[Dict]) -> Dict:
        """
        Analyze trends in the data
        """
        try:
            df = pd.DataFrame(data)
            
            # Determine value column
            value_col = df.select_dtypes(include=[np.number]).columns[0] if len(df.select_dtypes(include=[np.number]).columns) > 0 else None
            
            if value_col is None:
                raise ValueError("No numeric column found in data")
            
            values = df[value_col].values
            
            # Calculate trend
            x = np.arange(len(values))
            slope = np.polyfit(x, values, 1)[0]
            
            # Calculate percentage change
            if len(values) > 1:
                percentage_change = ((values[-1] - values[0]) / values[0]) * 100
            else:
                percentage_change = 0
            
            # Determine trend direction
            if slope > 0:
                trend_direction = "increasing"
            elif slope < 0:
                trend_direction = "decreasing"
            else:
                trend_direction = "stable"
            
            return {
                "trend_direction": trend_direction,
                "slope": float(slope),
                "percentage_change": float(percentage_change),
                "mean": float(np.mean(values)),
                "std": float(np.std(values)),
                "min": float(np.min(values)),
                "max": float(np.max(values))
            }
        except Exception as e:
            raise Exception(f"Trend analysis failed: {str(e)}")
    
    def detect_anomalies(self, data: List[Dict]) -> Dict:
        """
        Detect anomalies in the data using Isolation Forest
        """
        try:
            df = pd.DataFrame(data)
            
            # Determine value column
            value_col = df.select_dtypes(include=[np.number]).columns[0] if len(df.select_dtypes(include=[np.number]).columns) > 0 else None
            
            if value_col is None:
                raise ValueError("No numeric column found in data")
            
            values = df[value_col].values.reshape(-1, 1)
            
            # Scale values
            scaled_values = self.scaler.fit_transform(values)
            
            # Detect anomalies
            anomalies = self.isolation_forest.fit_predict(scaled_values)
            
            # Get anomaly indices
            anomaly_indices = np.where(anomalies == -1)[0].tolist()
            
            # Get anomaly values
            anomaly_values = []
            for idx in anomaly_indices:
                anomaly_values.append({
                    "index": int(idx),
                    "value": float(values[idx][0]),
                    "date": df.iloc[idx].get('date', df.iloc[idx].get('created_at', ''))
                })
            
            return {
                "total_anomalies": len(anomaly_indices),
                "anomaly_percentage": (len(anomaly_indices) / len(values)) * 100,
                "mean": float(np.mean(values)),
                "std_dev": float(np.std(values)),
                "anomalies": anomaly_values
            }
        except Exception as e:
            raise Exception(f"Anomaly detection failed: {str(e)}")
    
    def analyze_correlation(self, data: List[Dict]) -> Dict:
        """
        Analyze correlation between different metrics
        """
        try:
            df = pd.DataFrame(data)
            
            # Select numeric columns
            numeric_df = df.select_dtypes(include=[np.number])
            
            if len(numeric_df.columns) < 2:
                return {"message": "Insufficient numeric columns for correlation analysis"}
            
            # Calculate correlation matrix
            correlation_matrix = numeric_df.corr()
            
            # Convert to dictionary
            correlations = {}
            for col1 in correlation_matrix.columns:
                correlations[col1] = {}
                for col2 in correlation_matrix.columns:
                    correlations[col1][col2] = float(correlation_matrix.loc[col1, col2])
            
            return {
                "correlation_matrix": correlations,
                "strong_correlations": self._find_strong_correlations(correlation_matrix)
            }
        except Exception as e:
            raise Exception(f"Correlation analysis failed: {str(e)}")
    
    def _find_strong_correlations(self, correlation_matrix: pd.DataFrame, threshold: float = 0.7) -> List[Dict]:
        """
        Find strong correlations in the correlation matrix
        """
        strong_correlations = []
        for i in range(len(correlation_matrix.columns)):
            for j in range(i + 1, len(correlation_matrix.columns)):
                col1 = correlation_matrix.columns[i]
                col2 = correlation_matrix.columns[j]
                corr_value = correlation_matrix.loc[col1, col2]
                if abs(corr_value) >= threshold:
                    strong_correlations.append({
                        "metric1": col1,
                        "metric2": col2,
                        "correlation": float(corr_value)
                    })
        return strong_correlations

