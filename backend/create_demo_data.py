#!/usr/bin/env python
import os
import django
os.environ.setdefault('DJANGO_SETTINGS_MODULE', 'config.settings')
django.setup()

from locations.models import Category, Location, Season
from users.models import User

# Create categories
wine_category, _ = Category.objects.get_or_create(
    slug='winery',
    defaults={
        'name': 'Винодельня',
        'description': 'Винодельческие хозяйства и винотеки',
        'icon': 'wine_bar',
        'color': '#8B4513'
    }
)

nature_category, _ = Category.objects.get_or_create(
    slug='nature',
    defaults={
        'name': 'Природа',
        'description': 'Природные достопримечательности и парки',
        'icon': 'nature',
        'color': '#2E7D32'
    }
)

# Create seasons
for season_name, display_name in [('spring', 'Весна'), ('summer', 'Лето'), ('autumn', 'Осень'), ('winter', 'Зима')]:
    Season.objects.get_or_create(
        name=season_name,
        defaults={
            'description': f'{display_name} в Краснодарском крае',
            'start_month': [3, 6, 9, 12][['spring', 'summer', 'autumn', 'winter'].index(season_name)],
            'end_month': [5, 8, 11, 2][['spring', 'summer', 'autumn', 'winter'].index(season_name)],
        }
    )

# Create sample location
location, created = Location.objects.get_or_create(
    slug='abrau-durso',
    defaults={
        'name': 'Абрау-Дюрсо',
        'description': 'Известное винодельческое хозяйство на Черноморском побережье Кавказа, производящее шампанское и бренди в честь традиций русского шампанского.',
        'short_description': 'Популярное винодельческое хозяйство в Краснодарском крае',
        'category': wine_category,
        'latitude': '43.6953',
        'longitude': '39.3755',
        'address': 'Кабардинка, Краснодарский край',
        'image': None,
        'phone': '+7 (861) 350-01-20',
        'email': 'info@abrau-durso.ru',
        'website': 'https://www.abrau-durso.ru',
        'price_per_person': 500,
        'duration_hours': 3,
        'capacity': 100,
        'rating': 4.7,
        'reviews_count': 345,
        'featured': True,
        'verified': True,
    }
)

if created:
    location.best_season.set(Season.objects.filter(name__in=['spring', 'summer', 'autumn']))
    print(f"✓ Demo location '{location.name}' created successfully")
else:
    print(f"✓ Demo location '{location.name}' already exists")

# Create some test locations if needed
restaurant_category, _ = Category.objects.get_or_create(
    slug='restaurant',
    defaults={
        'name': 'Ресторан',
        'description': 'Рестораны и кафе',
        'icon': 'restaurant',
        'color': '#D84315'
    }
)

print("\n✓ Demo data created successfully!")
print("  - Categories: Винодельня, Природа, Ресторан")
print("  - Seasons: Весна, Лето, Осень, Зима")
print("  - Sample location: Абрау-Дюрсо")
