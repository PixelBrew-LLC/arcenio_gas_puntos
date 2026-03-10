import 'package:dio/dio.dart';
import '../constants.dart';
import '../storage/secure_storage_service.dart';

class DioClient {
  late final Dio dio;
  final SecureStorageService _storageService;

  DioClient(this._storageService) {
    dio = Dio(
      BaseOptions(
        baseUrl: AppConstants.baseUrl,
        connectTimeout: const Duration(seconds: 10),
        receiveTimeout: const Duration(seconds: 10),
        headers: {
          'Content-Type': 'application/json',
          'Accept': 'application/json',
        },
      ),
    );

    dio.interceptors.add(AuthInterceptor(_storageService));
  }
}

class AuthInterceptor extends Interceptor {
  final SecureStorageService _storageService;

  AuthInterceptor(this._storageService);

  @override
  void onRequest(
    RequestOptions options,
    RequestInterceptorHandler handler,
  ) async {
    // No agregar token a la ruta de login
    if (options.path.contains('/auth/login')) {
      return handler.next(options);
    }

    // Verificar que el token sea válido para hoy
    final isValid = await _storageService.isTokenValidForToday();
    if (!isValid) {
      // Token expirado por día calendario
      await _storageService.clearAll();
      return handler.reject(
        DioException(
          requestOptions: options,
          type: DioExceptionType.cancel,
          error: 'session_expired',
        ),
      );
    }

    final token = await _storageService.getToken();
    if (token != null) {
      options.headers['Authorization'] = 'Bearer $token';
    }

    return handler.next(options);
  }

  @override
  void onError(DioException err, ErrorInterceptorHandler handler) async {
    if (err.response?.statusCode == 401) {
      // Token rechazado por el servidor
      await _storageService.clearAll();
    }
    return handler.next(err);
  }
}
