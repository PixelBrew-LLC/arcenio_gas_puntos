import 'package:flutter/services.dart';

/// Utilidades para validar y formatear cédulas dominicanas.
/// Formato: ###-#######-# (11 dígitos separados 3-7-1)
class CedulaUtils {
  /// Elimina guiones de una cédula formateada → solo dígitos.
  static String unformat(String cedula) {
    return cedula.replaceAll('-', '').replaceAll(RegExp(r'[^0-9]'), '');
  }

  /// Formatea 11 dígitos como ###-#######-#
  static String format(String rawDigits) {
    final digits = unformat(rawDigits);
    if (digits.length != 11) return rawDigits;
    return '${digits.substring(0, 3)}-${digits.substring(3, 10)}-${digits.substring(10)}';
  }

  /// Valida la cédula dominicana usando el algoritmo de dígito verificador.
  /// Recibe dígitos sin formato (11 chars) o formateados.
  static bool isValid(String input) {
    final digits = unformat(input);
    if (digits.length != 11) return false;

    // Verificar que son solo números
    if (!RegExp(r'^\d{11}$').hasMatch(digits)) return false;

    // Algoritmo de validación dominicana
    final cedula10 = digits.substring(0, 10);
    final lastDigit = int.parse(digits.substring(10));

    const scale = [1, 2, 1, 2, 1, 2, 1, 2, 1, 2];
    int sum = 0;

    for (int i = 0; i < 10; i++) {
      final product = int.parse(cedula10[i]) * scale[i];
      // Sumar cada dígito del producto individualmente
      if (product >= 10) {
        sum += (product ~/ 10) + (product % 10);
      } else {
        sum += product;
      }
    }

    final nextTen = sum - (sum % 10) + 10;
    final expected = nextTen - sum;

    // Si expected es 10, el dígito verificador debe ser 0
    return (expected == 10 ? 0 : expected) == lastDigit;
  }

  /// Mensaje de error para validación (nullable = válido).
  static String? validate(String? value) {
    if (value == null || value.isEmpty) return 'Cédula requerida';
    final digits = unformat(value);
    if (digits.length != 11) return 'Debe tener 11 dígitos';
    if (!isValid(digits)) return 'Cédula inválida';
    return null;
  }
}

/// TextInputFormatter que aplica máscara ###-#######-#
class CedulaInputFormatter extends TextInputFormatter {
  @override
  TextEditingValue formatEditUpdate(
    TextEditingValue oldValue,
    TextEditingValue newValue,
  ) {
    // Extraer solo dígitos
    final digitsOnly = newValue.text.replaceAll(RegExp(r'[^0-9]'), '');

    // Limitar a 11 dígitos
    final limited = digitsOnly.length > 11
        ? digitsOnly.substring(0, 11)
        : digitsOnly;

    // Construir con guiones
    final buffer = StringBuffer();
    for (int i = 0; i < limited.length; i++) {
      if (i == 3 || i == 10) buffer.write('-');
      buffer.write(limited[i]);
    }

    final formatted = buffer.toString();

    // Calcular posición del cursor
    int cursorOffset = formatted.length;
    if (newValue.selection.baseOffset <= newValue.text.length) {
      // Contar dígitos hasta la posición original del cursor
      int digitCount = 0;
      for (
        int i = 0;
        i < newValue.selection.baseOffset && i < newValue.text.length;
        i++
      ) {
        if (RegExp(r'[0-9]').hasMatch(newValue.text[i])) {
          digitCount++;
        }
      }
      // Encontrar la posición en el texto formateado que corresponde a esa cantidad de dígitos
      int count = 0;
      cursorOffset = formatted.length;
      for (int i = 0; i < formatted.length; i++) {
        if (RegExp(r'[0-9]').hasMatch(formatted[i])) {
          count++;
          if (count == digitCount) {
            cursorOffset = i + 1;
            break;
          }
        }
      }
    }

    return TextEditingValue(
      text: formatted,
      selection: TextSelection.collapsed(
        offset: cursorOffset.clamp(0, formatted.length),
      ),
    );
  }
}
