import { isTauri } from '@tauri-apps/api/core';

/**
 * Rozpoznanie powłoki natywnej — jedno pytanie zadawane przez wszystkie mosty.
 *
 * Interfejs Danaco Console działa w dwóch miejscach naraz: w oknie powłoki
 * natywnej (`desktop/src-tauri/`) i w zwykłej przeglądarce. Każdy most do
 * powłoki musi najpierw ustalić, czy powłoka w ogóle jest — inaczej wywołanie
 * IPC rzuca wyjątkiem tam, gdzie żadnego IPC nie ma.
 *
 * Pytanie stoi w osobnym pliku, żeby most rdzenia nie zależał od mostu
 * katalogów tylko po to, by je zadać; oba pytają tak samo.
 *
 * Brak powłoki nie jest błędem (fail-open): `isTauri()` wykonane poza powłoką
 * ma prawo rzucić, a wtedy odpowiedzią jest zwyczajne „nie ma powłoki".
 */
export function czyPowlokaNatywna(): boolean {
  try {
    return isTauri();
  } catch {
    return false;
  }
}
