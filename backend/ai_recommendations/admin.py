"""Django admin configuration for AI recommendations app."""
from django.contrib import admin
from .models import RouteRecommendation, AIPromptLog


@admin.register(RouteRecommendation)
class RouteRecommendationAdmin(admin.ModelAdmin):
    list_display = ['user', 'title', 'season', 'score', 'user_feedback', 'created_at']
    list_filter = ['season', 'user_feedback', 'created_at']
    search_fields = ['title', 'user__username', 'description']
    readonly_fields = ['created_at']
    
    fieldsets = (
        ('Basic Info', {
            'fields': ('user', 'title', 'description')
        }),
        ('Route Data', {
            'fields': ('locations', 'estimated_duration', 'total_distance_km')
        }),
        ('AI Analysis', {
            'fields': ('ai_reasoning', 'score')
        }),
        ('User Interaction', {
            'fields': ('season', 'user_feedback', 'created_at')
        }),
    )


@admin.register(AIPromptLog)
class AIPromptLogAdmin(admin.ModelAdmin):
    list_display = ['user', 'success', 'execution_time_ms', 'tokens_used', 'created_at']
    list_filter = ['success', 'created_at']
    search_fields = ['user__username', 'user_query']
    readonly_fields = ['created_at', 'execution_time_ms', 'tokens_used']
    
    def has_add_permission(self, request):
        return False  # Prevent manual adding of logs
    
    def has_delete_permission(self, request, obj=None):
        return False  # Prevent deletion of logs
