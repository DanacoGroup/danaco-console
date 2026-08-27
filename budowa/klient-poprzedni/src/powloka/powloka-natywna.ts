import { isTauri } from '@tauri-apps/api/core';

/** Rozpoznanie powłoki natywnej — jedno pytanie zadawane przez wszystkie mosty do powłoki natywnej Tauri. */
export function czyPowlokaNatywna(): boolean {
  try {
    return isTauri();
  } catch {
    return false;
  }
}
