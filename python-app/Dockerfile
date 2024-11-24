# Use an official Python runtime as a parent image
FROM python:3.9-slim

# Set the working directory in the container
WORKDIR /app

# Copy requirements file and install dependencies
COPY requirements.txt .
RUN python3 -m pip install --no-cache-dir -r requirements.txt

# Copy the rest of the application code
COPY . .

# Run tests during the build process, running here because if test fails , no image is generated that is faulty.
RUN python3 -m pytest test-app.py  

# Expose the Flask app port
EXPOSE 5000

# Command to run the application
CMD ["python3", "app.py"]
