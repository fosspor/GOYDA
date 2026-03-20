"""Yandex AI service for generating route recommendations."""
import json
import logging
import time
from typing import Dict, List, Optional
import requests
from django.conf import settings

logger = logging.getLogger(__name__)


class YandexAIService:
    """Service for Yandex LLM API integration."""
    
    API_ENDPOINT = "https://llm.api.cloud.yandex.net/foundationModels/v1/completion"
    
    def __init__(self):
        self.folder_id = settings.YANDEX_FOLDER_ID
        self.api_key = settings.YANDEX_API_KEY
        self.iam_token = settings.YANDEX_IAM_TOKEN
    
    def _get_headers(self) -> Dict[str, str]:
        """Get request headers for Yandex API."""
        return {
            'Content-Type': 'application/json',
            'Authorization': f'Api-Key {self.api_key}' if self.api_key else f'Bearer {self.iam_token}',
            'x-folder-id': self.folder_id,
        }
    
    def generate_route_prompt(self, 
                            user_interests: List[str],
                            budget: int,
                            duration: int,
                            season: str,
                            group_info: Dict) -> str:
        """Generate a prompt for route recommendation."""
        
        prompt = f"""Вы - эксперт по туризму в Краснодарском крае.
        
Пожалуйста, создайте персональный маршрут с учетом следующих параметров:
- Интересы: {', '.join(user_interests)}
- Бюджет: {budget} рублей на противника
- Длительность: {duration} дней
- Сезон: {season}
- Группа: {json.dumps(group_info, ensure_ascii=False)}

Маршрут должен включать:
1. 3-5 интересных локаций Краснодарского края
2. Примерное время в пути между локациями
3. Рекомендуемое время посещения каждой локации
4. Специальные советы для данного сезона
5. Ориентировочные расходы

Ответ в формате JSON с ключами: locations, total_duration, tips, estimated_cost"""
        
        return prompt
    
    def call_completion_api(self, prompt: str, max_tokens: int = 2000) -> Optional[str]:
        """Call Yandex LLM API for text completion."""
        
        try:
            payload = {
                'modelUri': f'gpt://{self.folder_id}/yandexgpt-lite',
                'completionOptions': {
                    'stream': False,
                    'temperature': 0.7,
                    'maxTokens': max_tokens
                },
                'messages': [
                    {
                        'role': 'user',
                        'text': prompt
                    }
                ]
            }
            
            start_time = time.time()
            response = requests.post(
                self.API_ENDPOINT,
                headers=self._get_headers(),
                json=payload,
                timeout=30
            )
            execution_time = int((time.time() - start_time) * 1000)
            
            response.raise_for_status()
            result = response.json()
            
            # Extract completion
            if 'result' in result and 'alternatives' in result['result']:
                completion_text = result['result']['alternatives'][0]['message']['text']
                return completion_text
            
            logger.warning(f"Unexpected Yandex API response: {result}")
            return None
            
        except requests.RequestException as e:
            logger.error(f"Yandex API error: {e}")
            return None
    
    def generate_location_description(self, location_name: str, context: str) -> Optional[str]:
        """Generate compelling description for a location."""
        
        prompt = f"""Напишите увлекательное описание для туристического объекта:
        
Название: {location_name}
Контекст: {context}

Описание должно быть на русском языке, объемом 200-300 символов, интересным и информативным."""
        
        return self.call_completion_api(prompt, max_tokens=500)


# Singleton instance
_yandex_service = None

def get_yandex_service() -> YandexAIService:
    """Get or create Yandex AI service instance."""
    global _yandex_service
    if _yandex_service is None:
        _yandex_service = YandexAIService()
    return _yandex_service
