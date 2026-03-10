import 'package:dio/dio.dart';
import '../../../core/network/dio_client.dart';

class ClientRemoteDatasource {
  final DioClient _dioClient;

  ClientRemoteDatasource(this._dioClient);

  Future<Map<String, dynamic>> searchByCedula(String cedula) async {
    try {
      final response = await _dioClient.dio.get('/clients/$cedula');
      return response.data as Map<String, dynamic>;
    } on DioException catch (e) {
      if (e.response?.statusCode == 404) {
        throw Exception('Cliente no encontrado');
      }
      if (e.response != null) {
        throw Exception(e.response?.data['error'] ?? 'Error al buscar cliente');
      }
      throw Exception('Error de conexión con el servidor');
    }
  }

  Future<Map<String, dynamic>> createClient(Map<String, dynamic> data) async {
    try {
      final response = await _dioClient.dio.post('/clients', data: data);
      return response.data as Map<String, dynamic>;
    } on DioException catch (e) {
      if (e.response != null) {
        throw Exception(e.response?.data['error'] ?? 'Error al crear cliente');
      }
      throw Exception('Error de conexión con el servidor');
    }
  }

  Future<List<dynamic>> listClients() async {
    try {
      final response = await _dioClient.dio.get('/clients');
      return response.data as List<dynamic>;
    } on DioException catch (e) {
      if (e.response != null) {
        throw Exception(
          e.response?.data['error'] ?? 'Error al listar clientes',
        );
      }
      throw Exception('Error de conexión con el servidor');
    }
  }
}
