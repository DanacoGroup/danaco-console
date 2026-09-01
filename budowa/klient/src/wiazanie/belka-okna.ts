// Belka okna aplikacji wiązana z oknem powłoki: okno systemowe stoi bez ramy,
// więc zwinięcie, rozwinięcie, zamknięcie i przeciąganie niesie belka z warstwy
// projektowej — ta sama w oknie wejścia i w powłoce Centrum.

import { isTauri } from '@tauri-apps/api/core';
import { getCurrentWindow, type Window } from '@tauri-apps/api/window';

/** Czynności okna powłoki, które belka wywołuje; uprawnienia do nich wylicza `capabilities/domyslne.json`. */
type CzynnoscOkna = 'minimize' | 'toggleMaximize' | 'close';

/** Etykiety kontrolek belki: klucz warstwy projektowej → czynność okna. */
const CZYNNOSCI: ReadonlyArray<readonly [RegExp, CzynnoscOkna]> = [
  [/^(zwiń|zwin|minimalizuj)/i, 'minimize'],
  [/^(rozwiń|rozwin|maksymalizuj)/i, 'toggleMaximize'],
  [/^zamknij/i, 'close'],
];

/** Belki, które zastępują ramę okna: wejściowa przed uwierzytelnieniem i powłoki. */
const BELKI = '.dn-okno-wejsciowe-belka, .dn-belka';

/* Import biblioteki jest jawny, nie przez obiekt globalny `__TAURI__`: most
   globalny wystawiałby IPC każdemu skryptowi w dokumencie. Poza powłoką
   `getCurrentWindow` nie ma metadanych okna, więc pyta się o nie dopiero po
   rozpoznaniu powłoki. */
function oknoPowloki(): Window | null {
  return isTauri() ? getCurrentWindow() : null;
}

/** Wiąże belkę okna z oknem powłoki. Poza powłoką nie ma czego wiązać. */
export function zwiazBelkeOkna(): void {
  const okno = oknoPowloki();
  if (okno === null) return;

  document.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const przycisk = cel.closest('button[aria-label]');
    if (przycisk === null || przycisk.closest(BELKI) === null) return;
    const etykieta = przycisk.getAttribute('aria-label') ?? '';
    const czynnosc = CZYNNOSCI.find(([wzor]) => wzor.test(etykieta))?.[1];
    if (czynnosc === undefined) return;
    void okno[czynnosc]();
  });

  /* Przeciąganie idzie z belki, ale nie z kontrolek na niej — inaczej
     kliknięcie w „zamknij” zaczynałoby przesuwanie okna zamiast zamykać je. */
  document.addEventListener('pointerdown', (zdarzenie) => {
    if (zdarzenie.button !== 0) return;
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    if (cel.closest(BELKI) === null) return;
    if (cel.closest('button, a, input, [role="button"]') !== null) return;
    void okno.startDragging();
  });
}

/** Pokazuje okno powłoki, gdy strona ma już co pokazać. */
export function pokazOknoPowloki(): void {
  void oknoPowloki()?.show();
}
