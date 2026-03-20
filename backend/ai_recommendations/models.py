"""AI Recommendations service - Yandex LLM integration."""
from django.db import models
from django.contrib.auth import get_user_model

User = get_user_model()


class RouteRecommendation(models.Model):
    """AI-generated route recommendations."""
    
    user = models.ForeignKey(User, on_delete=models.CASCADE, related_name='ai_recommendations')
    title = models.CharField(max_length=200)
    description = models.TextField()
    
    # Route data
    locations = models.JSONField(help_text='List of location IDs and order')
    estimated_duration = models.IntegerField(help_text='Duration in hours')
    total_distance_km = models.FloatField(default=0)
    
    # Reasoning
    ai_reasoning = models.TextField(help_text='Why this route was recommended')
    score = models.FloatField(default=0.0, help_text='Recommendation score 0-1')
    
    # Metadata
    season = models.CharField(max_length=10)
    created_at = models.DateTimeField(auto_now_add=True)
    user_feedback = models.CharField(
        max_length=20,
        choices=[('liked', 'Liked'), ('neutral', 'Neutral'), ('disliked', 'Disliked')],
        null=True,
        blank=True
    )
    
    class Meta:
        ordering = ['-created_at']


class AIPromptLog(models.Model):
    """Logging of AI prompts for analysis and improvement."""
    
    user = models.ForeignKey(User, on_delete=models.CASCADE, related_name='ai_prompts')
    user_query = models.TextField()
    ai_response = models.TextField()
    used_details = models.JSONField(default=dict)
    execution_time_ms = models.IntegerField()
    tokens_used = models.IntegerField(null=True, blank=True)
    success = models.BooleanField(default=True)
    error_message = models.TextField(blank=True)
    created_at = models.DateTimeField(auto_now_add=True)
