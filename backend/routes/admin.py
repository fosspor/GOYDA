"""Django admin configuration for routes app."""
from django.contrib import admin
from .models import Route, RouteLocation, RouteInterest


class RouteLocationInline(admin.TabularInline):
    model = RouteLocation
    extra = 1
    fields = ['location', 'order', 'duration_hours', 'start_time', 'notes']


@admin.register(Route)
class RouteAdmin(admin.ModelAdmin):
    list_display = ['name', 'user', 'season', 'duration_days', 'is_public', 'featured']
    list_filter = ['season', 'is_public', 'featured', 'created_at']
    search_fields = ['name', 'user__username']
    readonly_fields = ['created_at', 'updated_at']
    inlines = [RouteLocationInline]
    
    fieldsets = (
        ('Basic Info', {
            'fields': ('name', 'slug', 'user')
        }),
        ('Route Details', {
            'fields': ('description', 'season', 'duration_days')
        }),
        ('Visibility', {
            'fields': ('is_public', 'featured')
        }),
        ('Timestamps', {
            'fields': ('created_at', 'updated_at'),
            'classes': ('collapse',)
        }),
    )


@admin.register(RouteLocation)
class RouteLocationAdmin(admin.ModelAdmin):
    list_display = ['route', 'location', 'order', 'duration_hours', 'start_time']
    list_filter = ['route', 'order']
    search_fields = ['route__name', 'location__name']


@admin.register(RouteInterest)
class RouteInterestAdmin(admin.ModelAdmin):
    list_display = ['user', 'interest', 'level']
    list_filter = ['interest', 'level']
    search_fields = ['user__username']
