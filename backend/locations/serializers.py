"""Locations serializers"""
from rest_framework import serializers
from .models import Location, Category

class CategorySerializer(serializers.ModelSerializer):
    class Meta:
        model = Category
        fields = ['id', 'name', 'slug', 'description', 'icon', 'color']


class LocationSerializer(serializers.ModelSerializer):
    category = CategorySerializer(read_only=True)
    
    class Meta:
        model = Location
        fields = [
            'id', 'name', 'slug', 'description', 'short_description',
            'category', 'address', 'image', 'rating', 'tags',
            'phone', 'email', 'website'
        ]
