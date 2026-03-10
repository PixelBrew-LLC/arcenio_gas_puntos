import 'package:dio/dio.dart';
import '../../../core/network/dio_client.dart';

class AuthRemoteDatasource {
  final DioClient _dioClient;

  AuthRemoteDatasource(this._dioClient);

  Future<Map<String, dynamic>> login(String username, String password) async {
    try {
      final response = await _dioClient.dio.post(
        '/auth/login',
        data: {'username': username, 'password': password},
      );
      return response.data as Map<String, dynamic>;
    } on DioException catch (e) {
      if (e.response != null) {
        throw Exception(e.response?.data['error'] ?? 'Error de autenticación');
      }
      throw Exception('Error de conexión con el servidor');
    }
  }

  Future<Map<String, dynamic>> getMe() async {
    try {
      final response = await _dioClient.dio.get('/auth/me');
      return response.data as Map<String, dynamic>;
    } on DioException catch (e) {
      if (e.response != null) {
        throw Exception(e.response?.data['error'] ?? 'Error al obtener datos');
      }
      throw Exception('Error de conexión con el servidor');
    }
  }
}
