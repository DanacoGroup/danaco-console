import type { Channel } from '../../../shared/contract';
import { opisOdmowyBledu } from '../komponenty/odmowa';
import type { Wynik } from '../protokol/kanal';
import type { FormularzKanalu, OdczytKanalu } from './formularz-kanalu';
import { niepowodzenie, potwierdzenie, type KomunikatZmiany } from './komunikat-zmiany';
import type { RejestrKanalow } from './rejestr-kanalow';
import type { WykazKanalow } from './wykaz-kanalow';
import type { ZrodloKanalow } from './zrodlo-kanalow';

/**
 * Trzy czynności panelu rejestru kanałów, każda jako funkcja nad wspólnym
 * kontekstem, nie jako gałąź jednego domknięcia. Kontekst wędruje jawnym
 * parametrem, więc żadna z tych funkcji nie zamyka się nad stanem panelu.
 *
 * Odmowa idzie dwiema drogami naraz: pasek komunikatów niesie kod oraz treść
 * rdzenia, a wykaz przechodzi w fazę `blad` z rolą `alert`. Wiersz przy odmowie
 * usunięcia zostaje na wykazie, bo panel nie kasuje niczego lokalnie — wykaz
 * odrysowuje się wyłącznie z rejestru.
 */
export interface KontekstKanalow {
  zrodlo: ZrodloKanalow;
  formularz: FormularzKanalu;
  wykaz: WykazKanalow;
  /** Rejestr tylko do odczytu; panel woła z niego wyłącznie `odswiez()`. */
  rejestr: RejestrKanalow;
  uzbrojenie: Uzbrojenie;
  zglos(komunikat: KomunikatZmiany): void;
}

/**
 * Uzbrojenie usuwania — dwa kroki zamiast jednego kliknięcia.
 *
 * To nie jest blokada: przycisk pozostaje czynny w każdej chwili, a pierwsze
 * naciśnięcie odpowiada — nazywa wiersz i mówi, że czynność jest nieodwracalna.
 * Blokadą byłoby `disabled`, którego tu nie ma.
 */
export interface Uzbrojenie {
  /** Kanał uzbrojony do usunięcia; pusty napis znaczy „rozbrojone". */
  kanal(): string;
  /** Uzbraja na wskazany wiersz i przepisuje etykietę przycisku. */
  uzbrojNa(idKanalu: string, ostrzezenie: string): void;
  /** Rozbraja i przywraca etykietę spoczynkową. */
  rozbroj(): void;
}

/** Etykieta spoczynkowa przycisku usuwania. */
const ETYKIETA_USUWANIA = 'Usuń wiersz rejestru';

export function utworzUzbrojenie(przycisk: HTMLButtonElement): Uzbrojenie {
  let uzbrojony = '';
  return {
    kanal: () => uzbrojony,
    uzbrojNa(idKanalu, ostrzezenie) {
      uzbrojony = idKanalu;
      przycisk.textContent = `Potwierdź usunięcie ${idKanalu}`;
      przycisk.setAttribute('aria-description', ostrzezenie);
    },
    rozbroj() {
      uzbrojony = '';
      przycisk.textContent = ETYKIETA_USUWANIA;
      przycisk.removeAttribute('aria-description');
    },
  };
}

export const OSTRZEZENIE_USUNIECIA =
  'Usunięcie wiersza jest NIEODWRACALNE. Kanał używany przez okno komunikacji rdzeń odrzuci ' +
  '(więź ON DELETE RESTRICT w schemacie bazy) — wiersz zostanie, a Ty zobaczysz odmowę z kodem.';

/** Wskazanie wiersza: formularz dostaje jego treść, uzbrojenie usuwania pada. */
export function wskazWiersz(k: KontekstKanalow, wskazany: Channel | undefined): void {
  k.formularz.wypelnij(wskazany);
  k.formularz.ostrzezRodzaj(wskazany !== undefined);
  k.uzbrojenie.rozbroj();
  k.zglos(
    wskazany === undefined
      ? potwierdzenie('Wskazanie porzucone — formularz zakłada nowy wiersz')
      : potwierdzenie(`Wskazano wiersz ${wskazany.id}`),
  );
}

/** `channel.add` — jedyna droga zakładania kanału w produkcie. */
export async function zalozKanal(k: KontekstKanalow): Promise<void> {
  const odczyt = k.formularz.odczytaj();
  if (odczyt.blad !== '') return k.zglos({ tresc: odczyt.blad, udany: false });
  if (odczyt.nazwa === '' || odczyt.rodzaj === '') {
    return k.zglos({
      tresc:
        'Założenie kanału: nazwa i rodzaj są w kontrakcie wymagane (ChannelAddRequest) — żądania nie wysłano.',
      udany: false,
    });
  }
  rozstrzygnij(k, 'Założenie kanału', await k.zrodlo.zaloz(zadanieZalozenia(odczyt)), (dodany) => {
    k.wykaz.wskaz(dodany.id);
    k.formularz.wypelnij(dodany);
    k.formularz.ostrzezRodzaj(true);
  });
}

/** `channel.update` — zmiana wiersza wskazanego na wykazie. */
export async function zapiszKanal(k: KontekstKanalow): Promise<void> {
  const wskazany = k.wykaz.wskazany();
  if (wskazany === undefined) {
    return k.zglos({ tresc: 'Zmiana kanału: wskaż najpierw wiersz wykazu.', udany: false });
  }
  const odczyt = k.formularz.odczytaj();
  if (odczyt.blad !== '') return k.zglos({ tresc: odczyt.blad, udany: false });
  rozstrzygnij(
    k,
    'Zmiana kanału',
    await k.zrodlo.zmien({ channelId: wskazany.id, ...zadanieZmiany(odczyt) }),
    (zmieniony) => k.formularz.wypelnij(zmieniony),
  );
}

/** `channel.remove` — dwa kroki, nigdy jeden. */
export async function usunKanal(k: KontekstKanalow): Promise<void> {
  const wskazany = k.wykaz.wskazany();
  if (wskazany === undefined) {
    return k.zglos({ tresc: 'Usunięcie kanału: wskaż najpierw wiersz wykazu.', udany: false });
  }
  if (k.uzbrojenie.kanal() !== wskazany.id) {
    k.uzbrojenie.uzbrojNa(wskazany.id, OSTRZEZENIE_USUNIECIA);
    return k.zglos({
      tresc: `${OSTRZEZENIE_USUNIECIA} Naciśnij ponownie, aby usunąć ${wskazany.id}.`,
      udany: false,
    });
  }
  k.uzbrojenie.rozbroj();
  rozstrzygnij(k, 'Usunięcie kanału', await k.zrodlo.usun(wskazany.id), () =>
    k.formularz.wypelnij(undefined),
  );
}

/**
 * `channel.check` — sprawdzenie, czy wskazany kanał odpowiada.
 *
 * Wynik NICZEGO nie warunkuje: nie wyłącza wiersza, nie wstrzymuje zapisu i nie
 * zmienia formularza. Jest odpowiedzią na pytanie Operatora „czy to działa",
 * a rozstrzygnięcie, co z nią zrobić, należy do niego.
 */
export async function sprawdzKanal(k: KontekstKanalow): Promise<void> {
  const wskazany = k.wykaz.wskazany();
  if (wskazany === undefined) {
    return k.zglos({ tresc: 'Sprawdzenie kanału: wskaż najpierw wiersz wykazu.', udany: false });
  }
  const wynik = await k.zrodlo.sprawdz(wskazany.id);
  if (!wynik.udany || wynik.wynik === undefined) {
    return k.zglos({
      tresc: `Sprawdzenie kanału: ${wynik.blad?.message ?? 'rdzeń odmówił bez powodu'}`,
      udany: false,
    });
  }
  const odpowiedz = wynik.wynik;
  if (!odpowiedz.reachable) {
    return k.zglos({
      tresc: `Kanał ${wskazany.id} nie odpowiedział: ${odpowiedz.detail ?? 'rdzeń nie podał powodu'}`,
      udany: false,
    });
  }
  k.zglos(
    potwierdzenie(
      `Kanał ${wskazany.id} odpowiedział` +
        (odpowiedz.latencyMs === undefined ? '' : ` w ${String(odpowiedz.latencyMs)} ms`),
    ),
  );
}

/**
 * `channel.credential.status` — STAN poświadczenia wskazanego kanału.
 *
 * Odpowiedź mówi, czy poświadczenie jest ustawione, jakiego jest rodzaju
 * i gdzie mieszka. Treści sekretu nie niesie i nie może nieść — panel nie ma
 * czego wyświetlić poza stanem.
 */
export async function stanPoswiadczeniaKanalu(k: KontekstKanalow): Promise<void> {
  const wskazany = k.wykaz.wskazany();
  if (wskazany === undefined) {
    return k.zglos({ tresc: 'Stan poświadczenia: wskaż najpierw wiersz wykazu.', udany: false });
  }
  const wynik = await k.zrodlo.stanPoswiadczenia(wskazany.id);
  if (!wynik.udany || wynik.wynik === undefined) {
    return k.zglos({
      tresc: `Stan poświadczenia: ${wynik.blad?.message ?? 'rdzeń odmówił bez powodu'}`,
      udany: false,
    });
  }
  const stan = wynik.wynik.status;
  const gdzie = stan.managedBy === undefined ? '' : ` (${stan.managedBy})`;
  k.zglos(
    potwierdzenie(
      stan.present
        ? `Poświadczenie kanału ${wskazany.id} jest ustawione${gdzie}.`
        : `Poświadczenie kanału ${wskazany.id} NIE jest ustawione${gdzie}.`,
    ),
  );
}

/**
 * Rozstrzygnięcie zapisu — jedno miejsce dla wszystkich trzech czynności.
 *
 * Po powodzeniu idzie `rejestr.odswiez()`: to zamówienie `channel.list`, które
 * rejestr rozgłasza wszystkim swoim czytelnikom. Panel nie dopisuje ani nie
 * kasuje niczego w rejestrze — nie ma tam drogi zapisu.
 */
export function rozstrzygnij<T>(
  k: KontekstKanalow,
  czynnosc: string,
  wynik: Wynik<T>,
  przyUdanym: (tresc: T) => void,
): void {
  if (!wynik.udany || wynik.wynik === undefined) {
    k.zglos(niepowodzenie(czynnosc, wynik.blad));
    k.wykaz.pokazOdmowe(opisOdmowyBledu(czynnosc, wynik.blad));
    return;
  }
  przyUdanym(wynik.wynik);
  k.rejestr.odswiez();
  k.wykaz.odrysuj();
  k.zglos(potwierdzenie(czynnosc));
}

/** Treść żądania `channel.add`; parametry idą tylko wtedy, gdy je podano. */
function zadanieZalozenia(odczyt: OdczytKanalu): {
  name: string;
  kind: string;
  model: string;
  enabled: boolean;
  config?: unknown;
} {
  return {
    name: odczyt.nazwa,
    kind: odczyt.rodzaj,
    model: odczyt.model,
    enabled: odczyt.czynny,
    ...(odczyt.parametry === undefined ? {} : { config: odczyt.parametry }),
  };
}

/** Treść żądania `channel.update` bez `channelId` — ten dokłada wywołujący. */
function zadanieZmiany(odczyt: OdczytKanalu): {
  name: string;
  model: string;
  enabled: boolean;
  config?: unknown;
} {
  return {
    name: odczyt.nazwa,
    model: odczyt.model,
    enabled: odczyt.czynny,
    ...(odczyt.parametry === undefined ? {} : { config: odczyt.parametry }),
  };
}
