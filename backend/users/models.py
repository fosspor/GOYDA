"""Users models - custom User model extending Django auth."""
from django.db import models
from django.contrib.auth.models import AbstractUser


class User(AbstractUser):
    """Custom User model with additional fields."""
    
    # Profile
    avatar = models.ImageField(upload_to='avatars/', blank=True, null=True)
    bio = models.TextField(blank=True)
    
    # Preferences
    preferred_regions = models.JSONField(default=list, blank=True)
    travel_style = models.CharField(
        max_length=50,
        choices=[
            ('budget', 'Budget'),
            ('comfort', 'Comfort'),
            ('luxury', 'Luxury'),
        ],
        default='comfort'
    )
    
    # Group indicators
    is_business = models.BooleanField(default=False, help_text='Business owner (farm, winery, etc.)')
    is_guide = models.BooleanField(default=False, help_text='Tour guide/operator')
    
    # Metadata
    date_joined = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)
    is_active = models.BooleanField(default=True)
    
    class Meta:
        verbose_name = 'User'
        verbose_name_plural = 'Users'
    
    def __str__(self):
        return self.username


class UserPreference(models.Model):
    """User preferences for route recommendations."""
    
    user = models.OneToOneField(User, on_delete=models.CASCADE, related_name='preference')
    min_budget = models.IntegerField(default=0)
    max_budget = models.IntegerField(default=100000)
    travel_duration_min = models.IntegerField(default=1, help_text='Minimum trip duration in days')
    travel_duration_max = models.IntegerField(default=30)
    group_size = models.IntegerField(default=1, help_text='Number of people in group')
    age_range = models.JSONField(default=dict, help_text='{"min": 0, "max": 100}')
    with_children = models.BooleanField(default=False)
    with_elderly = models.BooleanField(default=False)
    
    physical_activity_level = models.IntegerField(
        default=5, help_text='1-10 scale'
    )
    
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)
