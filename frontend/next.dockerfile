# Use an official Node.js runtime as a base image
FROM node:18

# Set the working directory in the Docker container
WORKDIR /app

# Set the environment variable
ENV NEXT_PUBLIC_API_URL=http://localhost:8080

# Copy package.json and package-lock.json into the Docker container
COPY ./package*.json ./

# Install the application's dependencies inside the Docker container
RUN npm install

# Copy the rest of the application code into the Docker container
COPY . .

# Build the Next.js application
RUN npm run build

# Make port 3000 available outside the Docker container
EXPOSE 3000

# Start the application
CMD ["npm", "run", "dev"]
