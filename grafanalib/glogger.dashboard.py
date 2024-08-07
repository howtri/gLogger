from grafanalib.core import (
    Dashboard, TimeSeries, Target, GridPos,
    OPS_FORMAT, SECONDS_FORMAT
)

# Define the dashboard
dashboard = Dashboard(
    title="gLogger Metrics Dashboard",
    description="Dashboard displaying custom gLogger metrics",
    tags=['glogger', 'metrics'],
    timezone="browser",
    panels=[
        # Panel for total gRPC requests
        TimeSeries(
            title="Total gRPC Requests",
            dataSource='Prometheus',
            targets=[
                Target(
                    expr='grpc_requests_total',
                    legendFormat="{{ method }}",
                    refId='A',
                ),
            ],
            unit=OPS_FORMAT,
            gridPos=GridPos(h=8, w=16, x=0, y=0),
        ),
        # Panel for gRPC request duration
        TimeSeries(
            title="gRPC Request Duration",
            dataSource='Prometheus',
            targets=[
                Target(
                    expr='grpc_request_duration_seconds',
                    legendFormat="{{ method }}",
                    refId='A',
                ),
            ],
            unit=SECONDS_FORMAT,
            gridPos=GridPos(h=8, w=16, x=0, y=10),
        ),
    ],
).auto_panel_ids()

# Save the generated dashboard JSON
if __name__ == "__main__":
    with open('glogger-dashboard.json', 'w') as f:
        f.write(dashboard.to_json_data())
