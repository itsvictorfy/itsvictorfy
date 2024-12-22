FROM nginx:alpine

# Remove default config if you want to supply your own
# RUN rm /etc/nginx/conf.d/default.conf

# Copy your custom nginx.conf
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Copy your static files into Nginx's html directory
COPY pages/. /usr/share/nginx/html/

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
