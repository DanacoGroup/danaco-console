import { utworzSekcjeBraku } from './sekcja-braku';
import type { SekcjaUstawien } from './sekcje';

/** Sekcja kont modeli jest pozycją odsyłającą do panelu kont w oknie modeli, nie drugą ramą do tych samych danych rdzenia. */
export function utworzSekcjeKontaModeli(): SekcjaUstawien {
  return utworzSekcjeBraku({
    wstep:
      'Konta modeli — klucze dostawców, konto domyślne i przypisania — prowadzi ' +
      'panel kont w oknie modeli. Tutaj ich nie ma świadomie: te same dane ' +
      'w dwóch oknach rozjechałyby się przy pierwszej zmianie.',
    pomiar: [
      'Rodzina `account.*` liczy pięć komend (add, list, remove, update, ' +
        'default.set) i wszystkie pięć jest obsłużonych w ' +
        '`budowa/client/src/modele/`.',
      'To samo rozstrzygnięcie zdjęło z tego rejestru pięć sekcji zbudowanych ' +
        'na `account.*` i `identity.*` — uzasadnienie stoi w nagłówku ' +
        '`sekcje.ts` i nie zostaje tu powtórzone.',
    ],
    odeslanie:
      'Droga: okno modeli → panel kont. Zmiana wykonana tam obowiązuje ' +
      'wszędzie — nie ma drugiego miejsca, w którym trzeba ją powtórzyć.',
  });
}
