import time
import docker  # Docker SDK for Python
import argparse  # For parsing command-line arguments

def scale_service(service_name, delta):
    # Create Docker client from environment configuration
    client = docker.from_env()
    try:
        # Get the specified service
        service = client.services.get(service_name)
    except docker.errors.NotFound:
        print(f"Service {service_name} not found!")
        return

    # Extract current replica count from service configuration
    current_replicas = service.attrs.get('Spec', {}) \
                                  .get('Mode', {}) \
                                  .get('Replicated', {}) \
                                  .get('Replicas', 0)
    # Calculate new replica count, ensuring it's not negative
    new_replicas = max(current_replicas + delta, 0)
    print(f"Scaling service '{service_name}' from {current_replicas} to {new_replicas}")

    # Update service's replica count using the scale method
    service.scale(new_replicas)

if __name__ == "__main__":
    # Parse command-line arguments 
    parser = argparse.ArgumentParser(description="Scale a Docker Swarm service after a delay.")
    parser.add_argument("service_name", type=str, help="Name of the Docker service")
    parser.add_argument("delta", type=int, help="Change in replica count (positive to scale up, negative to scale down)")
    parser.add_argument("--delay", type=int, default=30, help="Delay in seconds before scaling (default: 30)")
    args = parser.parse_args()

    # Delay execution for the specified number of seconds
    time.sleep(args.delay)
    scale_service(args.service_name, args.delta)

# Example usage:
# sudo python scale_service.py cs208_api-server-normal 1 --delay 30