# Use the official Node.js image
FROM node:18

# Set the working directory
WORKDIR /frontend

# Copy package.json and install dependencies
COPY package.json package-lock.json ./
RUN npm install

# Copy the rest of the application code
COPY . .

# Expose port 3000
EXPOSE 8080

# Start the React app
CMD ["npm", "start"]