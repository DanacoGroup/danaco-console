// Emiter rejestru komend klienta: wykłada pełny wykaz komend kontraktu jako
// literały Command.X w pliku interfejsu, żeby katalog operacji karmił się
// jednym źródłem prawdy zamiast dynamicznym przejściem po enumie.

import { linieNaglowka } from './naglowek.mjs';

/** Składa treść budowa/klient/src/wiazanie/rejestr-komend.ts z modelu kontraktu. */
export function emitujRejestrKlienta(model) {
  const naglowek = linieNaglowka(model.kontrakt, 'rejestr-komend.ts');
  const wpisy = model.komendy.map((k) => `  Command.${k.staly},`);
  return [
    '/*',
    ...naglowek.map((l) => (l ? ` * ${l}` : ' *')),
    ' */',
    '',
    "import { Command } from '../../../shared/contract.ts';",
    '',
    '/** Pełny wykaz komend kontraktu — źródło wykazu katalogu operacji interfejsu. */',
    'export const REJESTR_KOMEND: readonly Command[] = [',
    ...wpisy,
    '];',
    '',
  ].join('\n');
}
