import mysql.connector
from mysql.connector import Error
import os
from typing import List, Dict, Optional
from datetime import datetime
import json
import uuid

class Database:
    def __init__(self):
        self.connection = None
        self.connect()
    
    def connect(self):
        """Establish database connection"""
        if self.connection and self.connection.is_connected():
            return
        try:
            self.connection = mysql.connector.connect(
                host=os.getenv('DB_HOST', 'localhost'),
                port=int(os.getenv('DB_PORT', 3306)),
                user=os.getenv('DB_USER', 'root'),
                password=os.getenv('DB_PASSWORD', ''),
                database=os.getenv('DB_NAME', 'idash')
            )
        except Error as e:
            print(f"Error connecting to database: {e}")
            raise
    
    def _fetch_all(self, query: str, params: tuple) -> List[Dict]:
        try:
            self.connect()
            cursor = self.connection.cursor(dictionary=True)
            cursor.execute(query, params)
            results = cursor.fetchall()
            cursor.close()
            return results
        except Error as e:
            print(f"Database query failed: {e}")
            return []
    
    def get_sales_data(self, start_date: str, end_date: str, company_id: Optional[int] = None,
                       channel_id: Optional[int] = None, product_id: Optional[int] = None) -> List[Dict]:
        """Fetch historical sales data over a window"""
        query = """
            SELECT data_month, channel_id, channel_name, billed, primary_collection, ims,
                   lifting_target, ims_target
            FROM channelwise_monthly_report
            WHERE data_month BETWEEN %s AND %s
        """
        params: List = [start_date, end_date]
        if channel_id:
            query += " AND channel_id = %s"
            params.append(channel_id)
        query += " ORDER BY data_month ASC"
        return self._fetch_all(query, tuple(params))
    
    def get_sales_breakdown(self, start_date: str, end_date: str, granularity: str = "channel") -> List[Dict]:
        """Aggregate sales history by channel or product."""
        if granularity == "product":
            query = """
                SELECT year_month AS period, product_id, product_name, channel_id,
                       SUM(value) AS sales_value, SUM(qty) AS quantity
                FROM best_selling_products
                WHERE year_month BETWEEN %s AND %s
                GROUP BY year_month, product_id, product_name, channel_id
                ORDER BY year_month ASC
            """
        else:
            query = """
                SELECT DATE_FORMAT(data_month, '%Y-%m-01') AS period,
                       channel_id,
                       channel_name,
                       SUM(billed) AS sales_value,
                       SUM(primary_collection) AS collection_value,
                       SUM(lifting_target) AS target_value
                FROM channelwise_monthly_report
                WHERE data_month BETWEEN %s AND %s
                GROUP BY period, channel_id, channel_name
                ORDER BY period ASC
            """
        return self._fetch_all(query, (start_date, end_date))
    
    def get_inventory_data(self, start_date: str, end_date: str, product_id: Optional[int] = None) -> List[Dict]:
        query = """
            SELECT month, company_id, gl_id, amount
            FROM inventory_raw_datas
            WHERE month BETWEEN %s AND %s
        """
        params: List = [start_date, end_date]
        if product_id:
            query += " AND gl_id = %s"
            params.append(product_id)
        query += " ORDER BY month ASC"
        return self._fetch_all(query, tuple(params))
    
    def get_current_inventory_snapshot(self) -> List[Dict]:
        query = """
            SELECT gl_id, SUM(amount) AS on_hand_value, COUNT(*) AS observations
            FROM inventory_raw_datas
            GROUP BY gl_id
        """
        return self._fetch_all(query, tuple())
    
    def get_financial_data(self, start_date: str, end_date: str, budget_category: Optional[str] = None) -> List[Dict]:
        query = """
            SELECT month, category_id, department_id, budget_amount, actual_amount
            FROM budget_summaries
            WHERE month BETWEEN %s AND %s
        """
        params: List = [start_date, end_date]
        if budget_category:
            query += " AND category_id = %s"
            params.append(budget_category)
        query += " ORDER BY month ASC"
        return self._fetch_all(query, tuple(params))
    
    def get_expense_breakdown(self, start_date: str, end_date: str) -> List[Dict]:
        query = """
            SELECT month, expense_id, budget_amount, actual_amount
            FROM budget_monthlies
            WHERE month BETWEEN %s AND %s
            ORDER BY month ASC
        """
        return self._fetch_all(query, (start_date, end_date))
    
    def get_production_data(self, start_date: str, end_date: str, factory_id: Optional[int] = None,
                            product_id: Optional[int] = None) -> List[Dict]:
        query = """
            SELECT production_date, factory_id, planned_output, actual_output
            FROM production_efficiency
            WHERE production_date BETWEEN %s AND %s
        """
        params: List = [start_date, end_date]
        if factory_id:
            query += " AND factory_id = %s"
            params.append(factory_id)
        query += " ORDER BY production_date ASC"
        return self._fetch_all(query, tuple(params))
    
    def get_metric_series(self, table: str, date_column: str, value_column: str,
                           start_date: str, end_date: str) -> List[Dict]:
        query = f"""
            SELECT {date_column} AS metric_date, {value_column} AS metric_value
            FROM {table}
            WHERE {date_column} BETWEEN %s AND %s
            ORDER BY {date_column} ASC
        """
        return self._fetch_all(query, (start_date, end_date))
    
    def save_forecast(self, forecast: Dict) -> str:
        """Return a generated forecast identifier (persistence handled upstream)."""
        forecast_id = forecast.get("forecast_id") or f"fc_{uuid.uuid4().hex[:12]}"
        return forecast_id
    
    def close(self):
        """Close database connection"""
        if self.connection and self.connection.is_connected():
            self.connection.close()

