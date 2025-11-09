import InsightListCard from '../Dashboard/InsightListCard'

const BIInsightsPanel = ({ insights = [] }) => (
  <InsightListCard title="AI Insights" insights={insights} iconColor="warning" dense />
)

export default BIInsightsPanel
