"""Django admin configuration for locations app."""
from django.contrib import admin
from django.utils.html import format_html
from .models import Category, Season, Location, LocationReview, LocationPhoto


@admin.register(Category)
class CategoryAdmin(admin.ModelAdmin):
    list_display = ['name', 'slug', 'color_preview']
    list_editable = ['slug']
    search_fields = ['name']
    
    def color_preview(self, obj):
        return format_html(
            '<div style="width: 30px; height: 30px; background-color: {0}; border-radius: 4px;"></div>',
            obj.color
        )
    color_preview.short_description = 'Color'


@admin.register(Season)
class SeasonAdmin(admin.ModelAdmin):
    list_display = ['name', 'start_month', 'end_month']
    readonly_fields = ['description']


@admin.register(Location)
class LocationAdmin(admin.ModelAdmin):
    list_display = ['name', 'category', 'address', 'rating', 'featured', 'verified']
    list_filter = ['category', 'featured', 'verified', 'best_season']
    search_fields = ['name', 'address']
    readonly_fields = ['created_at', 'updated_at', 'location_preview']
    
    fieldsets = (
        ('Basic Info', {
            'fields': ('name', 'slug', 'description', 'short_description')
        }),
        ('Location', {
            'fields': ('category', 'location', 'address', 'region', 'location_preview')
        }),
        ('Contact', {
            'fields': ('phone', 'email', 'website', 'instagram')
        }),
        ('Media', {
            'fields': ('image', 'thumbnail', 'video_url', 'gallery_images')
        }),
        ('Details', {
            'fields': ('price_per_person', 'duration_hours', 'capacity', 'amenities', 'accessibility')
        }),
        ('Seasonal & Rating', {
            'fields': ('best_season', 'weather_dependent', 'rating', 'reviews_count')
        }),
        ('Tags & Metadata', {
            'fields': ('tags', 'featured', 'verified')
        }),
        ('Timestamps', {
            'fields': ('created_at', 'updated_at'),
            'classes': ('collapse',)
        }),
    )
    
    def location_preview(self, obj):
        if obj.location:
            return format_html(
                '<p><strong>Coordinates:</strong> {0}, {1}</p>',
                obj.location.y,  # latitude
                obj.location.x    # longitude
            )
        return 'No location set'
    location_preview.short_description = 'Geographic Preview'


@admin.register(LocationReview)
class LocationReviewAdmin(admin.ModelAdmin):
    list_display = ['location', 'user', 'rating', 'created_at']
    list_filter = ['rating', 'created_at', 'location']
    search_fields = ['location__name', 'user__username']
    readonly_fields = ['created_at', 'updated_at']


@admin.register(LocationPhoto)
class LocationPhotoAdmin(admin.ModelAdmin):
    list_display = ['location', 'user', 'created_at', 'image_preview']
    list_filter = ['created_at', 'location']
    search_fields = ['location__name', 'user__username']
    readonly_fields = ['created_at', 'image_preview']
    
    def image_preview(self, obj):
        if obj.image:
            return format_html(
                '<img src="{0}" width="100" height="100" style="object-fit: cover; border-radius: 4px;"/>',
                obj.image.url
            )
        return 'No image'
    image_preview.short_description = 'Preview'
