"""
Models for locations app - storing attractions, restaurants, wineries, etc.
"""

from django.db import models
from django.core.validators import MinValueValidator, MaxValueValidator
from django.utils import timezone


class Category(models.Model):
    """Categories of locations: restaurants, wineries, nature, etc."""
    
    CATEGORY_CHOICES = [
        ('winery', 'Винодельня'),
        ('farm', 'Ферма'),
        ('restaurant', 'Ресторан'),
        ('guesthouse', 'Гостевой дом'),
        ('nature', 'Природа'),
        ('history', 'История'),
        ('arts', 'Искусство'),
        ('adventure', 'Активный отдых'),
        ('health', 'Оздоровление'),
        ('culture', 'Культура'),
        ('market', 'Рынок/Лавка'),
        ('workshop', 'Мастерская'),
    ]
    
    name = models.CharField(max_length=100)
    slug = models.SlugField(unique=True)
    description = models.TextField(blank=True)
    icon = models.CharField(max_length=50, default='location_on')
    color = models.CharField(max_length=7, default='#FF6B6B')
    
    class Meta:
        verbose_name_plural = 'Categories'
        ordering = ['name']
    
    def __str__(self):
        return self.name


class Season(models.Model):
    """Seasonal availability and recommendations."""
    
    SEASON_CHOICES = [
        ('spring', 'Весна'),
        ('summer', 'Лето'),
        ('autumn', 'Осень'),
        ('winter', 'Зима'),
    ]
    
    name = models.CharField(max_length=20, choices=SEASON_CHOICES, unique=True)
    description = models.TextField()
    start_month = models.IntegerField()
    end_month = models.IntegerField()
    
    def __str__(self):
        return self.get_name_display()


class Location(models.Model):
    """Main location model for attractions, restaurants, etc."""
    
    # Basic info
    name = models.CharField(max_length=200)
    slug = models.SlugField(unique=True)
    description = models.TextField()
    short_description = models.CharField(max_length=500)
    
    # Location
    category = models.ForeignKey(Category, on_delete=models.CASCADE, related_name='locations')
    latitude = models.DecimalField(max_digits=9, decimal_places=6, validators=[MinValueValidator(-90), MaxValueValidator(90)])
    longitude = models.DecimalField(max_digits=9, decimal_places=6, validators=[MinValueValidator(-180), MaxValueValidator(180)])
    address = models.CharField(max_length=300)
    region = models.CharField(max_length=100, default='Краснодарский край')
    
    # Media
    image = models.ImageField(upload_to='locations/', blank=True, null=True)
    thumbnail = models.ImageField(upload_to='locations/thumbnails/', blank=True, null=True)
    gallery_images = models.JSONField(default=list, blank=True, help_text='List of gallery image URLs')
    video_url = models.URLField(blank=True, null=True, help_text='URL to 360/video tour')
    
    # Contact
    phone = models.CharField(max_length=20, blank=True)
    email = models.EmailField(blank=True)
    website = models.URLField(blank=True)
    instagram = models.CharField(max_length=100, blank=True)
    
    # Details
    working_hours = models.JSONField(default=dict, blank=True, help_text='Days and hours of operation')
    price_per_person = models.DecimalField(max_digits=10, decimal_places=2, null=True, blank=True)
    duration_hours = models.IntegerField(null=True, blank=True, help_text='Duration of visit in hours')
    capacity = models.IntegerField(null=True, blank=True, help_text='Max visitors at once')
    
    # Amenities
    amenities = models.JSONField(default=list, blank=True, help_text='List of amenities: wifi, parking, etc.')
    accessibility = models.JSONField(default=list, blank=True, help_text='Accessibility features')
    
    # Seasonal
    best_season = models.ManyToManyField(Season, related_name='locations')
    weather_dependent = models.BooleanField(default=False)
    
    # Rating & Reviews
    rating = models.DecimalField(
        max_digits=3,
        decimal_places=2,
        default=0.0,
        validators=[MinValueValidator(0), MaxValueValidator(5)]
    )
    reviews_count = models.IntegerField(default=0)
    
    # Tags for search
    tags = models.JSONField(default=list, blank=True, help_text='Search tags: organic, vegan, family-friendly, etc.')
    
    # Metadata
    featured = models.BooleanField(default=False)
    verified = models.BooleanField(default=False)
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)
    
    objects = models.Manager()
    
    class Meta:
        ordering = ['-featured', '-rating']
        indexes = [
            models.Index(fields=['category', '-rating']),
            models.Index(fields=['featured']),
        ]
    
    def __str__(self):
        return self.name
    
    def to_dict(self):
        """Serialize for API responses."""
        return {
            'id': self.id,
            'name': self.name,
            'slug': self.slug,
            'category': self.category.slug,
            'description': self.description,
            'address': self.address,
            'latitude': self.location.y,
            'longitude': self.location.x,
            'image': self.image.url if self.image else None,
            'rating': float(self.rating),
            'tags': self.tags,
            'best_season': [s.name for s in self.best_season.all()],
        }


class LocationReview(models.Model):
    """User reviews for locations."""
    
    location = models.ForeignKey(Location, on_delete=models.CASCADE, related_name='user_reviews')
    user = models.ForeignKey('users.User', on_delete=models.CASCADE, related_name='location_reviews')
    rating = models.IntegerField(validators=[MinValueValidator(1), MaxValueValidator(5)])
    title = models.CharField(max_length=100)
    text = models.TextField()
    visited_date = models.DateField()
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)
    
    class Meta:
        unique_together = ['location', 'user']
        ordering = ['-created_at']
    
    def __str__(self):
        return f"{self.location.name} - {self.rating}/5"


class LocationPhoto(models.Model):
    """User photos for locations."""
    
    location = models.ForeignKey(Location, on_delete=models.CASCADE, related_name='photos')
    user = models.ForeignKey('users.User', on_delete=models.CASCADE, related_name='location_photos')
    image = models.ImageField(upload_to='location_photos/')
    caption = models.CharField(max_length=500, blank=True)
    created_at = models.DateTimeField(auto_now_add=True)
    
    class Meta:
        ordering = ['-created_at']
