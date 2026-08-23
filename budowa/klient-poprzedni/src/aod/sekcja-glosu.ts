import { opisOdmowyAod } from './odmowy-aod';
import { utworzAkapit, utworzPodtytul } from './pola-wykazu';
import type { ZrodloAod } from './zrodlo-komend';

/**
 * Sekcja wydająca polecenie głosowe do asystenta (`aod.voice.command`).
 *
 * Polecenie wydaje się samą transkrypcją — pola `audioRef` i `transcript`
 * kontraktu są opcjonalne, a nakładka nagrań nie tworzy, więc przycisku
 * mikrofonu tu nie ma. Sekcja rozmowy jest osobną drogą: `aod.chat.send` niesie
 * zdanie do okna rozmowy, `aod.voice.command` kieruje polecenie do asystenta
 * i ma własne pole `speak` na odpowiedź syntezą mowy.
 */
export interface SekcjaGlosu {
  element: HTMLElement;
  /** Wskazuje sesję, do której kierowane jest polecenie bez wskazanego okna. */
  ustawSesje(sessionId?: string): void;
}

export interface OpisSekcjiGlosu {
  zrodlo: ZrodloAod;
  /** Krótkie potwierdzenie czynności na pasku okna. */
  zamelduj(zdanie: string, udane: boolean): void;
}

export function utworzSekcjeGlosu(opis: OpisSekcjiGlosu): SekcjaGlosu {
  let sesja: string | undefined;

  const tresc = document.createElement('textarea');
  tresc.className = 'ao-pole ao-pole--tresc';
  tresc.rows = 2;
  tresc.placeholder = 'treść polecenia dla asystenta';
  tresc.setAttribute('aria-label', 'Transkrypcja polecenia głosowego');

  const okno = document.createElement('input');
  okno.type = 'text';
  okno.className = 'ao-pole';
  okno.placeholder = 'okno asystenta — puste kieruje do ogniskowanego';
  okno.setAttribute('aria-label', 'Okno asystenta wykonujące polecenie');

  const mowa = document.createElement('input');
  mowa.type = 'checkbox';
  mowa.className = 'dn-przelacznik';
  mowa.id = 'ao-glos-synteza';

  const etykietaMowy = document.createElement('label');
  etykietaMowy.className = 'ao-etykieta';
  etykietaMowy.htmlFor = mowa.id;
  etykietaMowy.textContent = 'Odczytaj odpowiedź syntezą mowy';

  const wyslij = document.createElement('button');
  wyslij.type = 'button';
  wyslij.className = 'dn-btn dn-btn--sm dn-btn--atrament';
  wyslij.textContent = 'Wydaj polecenie';

  const meldunek = document.createElement('div');
  meldunek.className = 'ao-meldunek';

  const granica = utworzAkapit(
    'ao-granica',
    'Nagrania nakładka nie robi — nagrywanie i rozpoznawanie mowy należą do silnika ' +
      'mowy. Polecenie wydajesz treścią; pole audioRef kontraktu zostaje puste.',
  );

  const element = document.createElement('section');
  element.className = 'ao-sekcja ao-sekcja--glos';
  element.append(
    utworzPodtytul('Polecenie głosowe do asystenta (aod.voice.command)'),
    tresc,
    okno,
    mowa,
    etykietaMowy,
    wyslij,
    meldunek,
    granica,
  );

  async function wydaj(): Promise<void> {
    const wpisana = tresc.value.trim();
    const wskazaneOkno = okno.value.trim();
    // Przycisk pozostaje czynny także przy pustym polu: pustą treść ocenia
    // rdzeń i to on zwraca odmowę.
    const wynik = await opis.zrodlo.wydajPolecenieGlosowe({
      ...(wpisana === '' ? {} : { transcript: wpisana }),
      ...(wskazaneOkno === '' ? {} : { windowId: wskazaneOkno }),
      ...(wskazaneOkno === '' && sesja !== undefined ? { sessionId: sesja } : {}),
      speak: mowa.checked,
    });

    if (!wynik.udany || wynik.wynik === undefined) {
      meldunek.replaceChildren(
        utworzAkapit('ao-odmowa', opisOdmowyAod('Wydanie polecenia głosowego', 'glos', wynik.blad)),
      );
      opis.zamelduj('Polecenie głosowe odrzucone.', false);
      return;
    }

    // Meldunek pokazuje transkrypcję zwróconą przez rdzeń, nie wpisaną —
    // rdzeń mógł ją poprawić.
    const oddana = wynik.wynik.transcript;
    meldunek.replaceChildren(
      utworzAkapit('ao-pokwitowanie', `Asystent przyjął polecenie: „${oddana}".`),
    );
    tresc.value = '';
    opis.zamelduj('Polecenie głosowe wydane.', true);
  }

  wyslij.addEventListener('click', () => void wydaj());

  return {
    element,
    ustawSesje(sessionId) {
      sesja = sessionId;
    },
  };
}
