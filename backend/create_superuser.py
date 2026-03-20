#!/usr/bin/env python
import os
import django
os.environ.setdefault('DJANGO_SETTINGS_MODULE', 'config.settings')
django.setup()

from django.contrib.auth import get_user_model

User = get_user_model()

# Create superuser
username = 'admin'
email = 'admin@krasnodar-tourism.local'
password = 'admin123456'

if not User.objects.filter(username=username).exists():
    User.objects.create_superuser(username, email, password)
    print(f"✓ Superuser '{username}' created successfully")
    print(f"  Email: {email}")
    print(f"  Password: {password}")
else:
    print(f"✓ Superuser '{username}' already exists")
