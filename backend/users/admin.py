"""Django admin configuration for users app."""
from django.contrib import admin
from django.contrib.auth.admin import UserAdmin as BaseUserAdmin
from .models import User, UserPreference


class UserPreferenceInline(admin.StackedInline):
    model = UserPreference
    extra = 0
    fields = [
        'min_budget', 'max_budget', 'travel_duration_min', 'travel_duration_max',
        'group_size', 'with_children', 'with_elderly', 'physical_activity_level'
    ]


@admin.register(User)
class UserAdmin(BaseUserAdmin):
    list_display = ['username', 'email', 'is_business', 'is_guide', 'date_joined']
    list_filter = ['is_business', 'is_guide', 'is_active', 'date_joined']
    search_fields = ['username', 'email']
    inlines = [UserPreferenceInline]
    
    fieldsets = BaseUserAdmin.fieldsets + (
        ('Profile', {
            'fields': ('avatar', 'bio')
        }),
        ('Preferences', {
            'fields': ('preferred_regions', 'travel_style')
        }),
        ('User Type', {
            'fields': ('is_business', 'is_guide')
        }),
    )


@admin.register(UserPreference)
class UserPreferenceAdmin(admin.ModelAdmin):
    list_display = ['user', 'min_budget', 'max_budget', 'group_size', 'physical_activity_level']
    list_filter = ['with_children', 'with_elderly']
    search_fields = ['user__username']
