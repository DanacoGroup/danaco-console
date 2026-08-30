// Belka okna aplikacji wiązana z oknem powłoki: okno systemowe stoi bez ramy,
// więc zwinięcie, rozwinięcie, zamknięcie i przeciąganie niesie belka z warstwy
// projektowej — ta sama w oknie wejścia i w powłoce Centrum.

type OknoPowloki = {
  minimize: () => Promise<void>;
  toggleMaximize: () => Promise<void>;
  close: () => Promise<void>;
  startDragging: () => Promise<void>;
  show: () => Promise<void>;
};

type MostTauri = { window?: { getCurrentWindow?: () => OknoPowloki } };

/** Etykiety kontrolek belki: klucz warstwy projektowej → czynność okna. */
const CZYNNOSCI: ReadonlyArray<readonly [RegExp, keyof OknoPowloki]> = [
  [/^(zwiń|zwin|minimalizuj)/i, 'minimize'],
  [/^(rozwiń|rozwin|maksymalizuj)/i, 'toggleMaximize'],
  [/^zamknij/i, 'close'],
];

/** Belki, które zastępują ramę okna: wejściowa przed uwierzytelnieniem i powłoki. */
const BELKI = '.dn-okno-wejsciowe-belka, .dn-belka';

function oknoPowloki(): OknoPowloki | null {
  const most = (globalThis as { __TAURI__?: MostTauri }).__TAURI__;
  const pobierz = most?.window?.getCurrentWindow;
  return pobierz === undefined ? null : pobierz();
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
