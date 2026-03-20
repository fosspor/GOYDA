"""Routes models"""
from django.db import models
from django.contrib.auth import get_user_model
from locations.models import Location

User = get_user_model()


class Route(models.Model):
    """Personalized travel routes."""
    
    SEASON_CHOICES = [
        ('spring', 'Весна'),
        ('summer', 'Лето'),
        ('autumn', 'Осень'),
        ('winter', 'Зима'),
    ]
    
    name = models.CharField(max_length=200)
    slug = models.SlugField(unique=True)
    description = models.TextField()
    
    # Creator & Ownership
    user = models.ForeignKey(User, on_delete=models.CASCADE, related_name='routes')
    is_public = models.BooleanField(default=False)
    
    # Content
    locations = models.ManyToManyField(Location, through='RouteLocation', related_name='routes')
    season = models.CharField(max_length=10, choices=SEASON_CHOICES)
    duration_days = models.IntegerField()
    
    # Metadata
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)
    featured = models.BooleanField(default=False)
    
    class Meta:
        ordering = ['-created_at']
    
    def __str__(self):
        return self.name


class RouteLocation(models.Model):
    """Order and details of locations in a route."""
    
    route = models.ForeignKey(Route, on_delete=models.CASCADE)
    location = models.ForeignKey(Location, on_delete=models.CASCADE)
    order = models.IntegerField()
    notes = models.TextField(blank=True)
    duration_hours = models.IntegerField(null=True, blank=True)
    start_time = models.TimeField(null=True, blank=True)
    
    class Meta:
        unique_together = ['route', 'location']
        ordering = ['order']


class RouteInterest(models.Model):
    """User interests for personalization."""
    
    INTERESTS = [
        ('gastronomy', 'Гастрономия'),
        ('wine', 'Вино'),
        ('nature', 'Природа'),
        ('history', 'История'),
        ('art', 'Искусство'),
        ('adventure', 'Приключения'),
        ('wellness', 'Оздоровление'),
        ('photography', 'Фотография'),
        ('family', 'Семейный отдых'),
        ('culture', 'Культура'),
    ]
    
    user = models.ForeignKey(User, on_delete=models.CASCADE, related_name='interests')
    interest = models.CharField(max_length=50, choices=INTERESTS)
    level = models.IntegerField(default=5, help_text='1-10 preference level')
    
    class Meta:
        unique_together = ['user', 'interest']
