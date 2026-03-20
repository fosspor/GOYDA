"""Locations views - API endpoints for attractions and restaurants"""
from rest_framework import viewsets, filters, status
from rest_framework.decorators import action
from rest_framework.response import Response
from rest_framework.pagination import PageNumberPagination
from django_filters.rest_framework import DjangoFilterBackend
from .models import Location, Category, Season
from .serializers import LocationSerializer, CategorySerializer


class LocationPagination(PageNumberPagination):
    """Pagination for location lists"""
    page_size = 12
    page_size_query_param = 'page_size'
    max_page_size = 100


class LocationViewSet(viewsets.ModelViewSet):
    """
    API endpoint for locations (attractions, restaurants, wineries, etc.)
    
    Supports:
    - GET /api/locations/ - List all locations with pagination
    - GET /api/locations/?search=wine - Search locations
    - GET /api/locations/?category=1 - Filter by category
    - GET /api/locations/?featured=true - Get featured only
    - GET /api/locations/featured/ - Get featured locations
    - GET /api/locations/{id}/nearby/ - Get nearby locations
    """
    queryset = Location.objects.all()
    serializer_class = LocationSerializer
    pagination_class = LocationPagination
    filter_backends = [DjangoFilterBackend, filters.SearchFilter, filters.OrderingFilter]
    filterset_fields = ['category', 'featured', 'verified']
    search_fields = ['name', 'description', 'short_description', 'tags']
    ordering_fields = ['rating', 'created_at', 'reviews_count']
    ordering = ['-featured', '-rating']
    
    @action(detail=False, methods=['get'])
    def featured(self, request):
        """Get featured locations"""
        featured_locations = self.get_queryset().filter(featured=True)
        page = self.paginate_queryset(featured_locations)
        if page is not None:
            serializer = self.get_serializer(page, many=True)
            return self.get_paginated_response(serializer.data)
        serializer = self.get_serializer(featured_locations, many=True)
        return Response(serializer.data)
    
    @action(detail=True, methods=['get'])
    def nearby(self, request, pk=None):
        """Get nearby locations (same category)"""
        location = self.get_object()
        nearby = Location.objects.filter(
            category=location.category
        ).exclude(pk=location.pk)[:5]
        serializer = self.get_serializer(nearby, many=True)
        return Response(serializer.data)


class CategoryViewSet(viewsets.ReadOnlyModelViewSet):
    """API endpoint for location categories"""
    queryset = Category.objects.all()
    serializer_class = CategorySerializer
    filter_backends = [filters.SearchFilter]
    search_fields = ['name', 'description']
