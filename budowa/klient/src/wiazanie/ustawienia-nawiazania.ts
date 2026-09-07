// Mobile i Always On Display w sekcji Nawiązań: stan platformy widziany
// z telefonu, sterowanie procesami oraz nakładka towarzysząca.
import {
  AodEventClass,
  AodMuteKind,
  Command,
  MobileProcessControl,
} from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { oglos } from './ogloszenie.ts';
import { zapytajWOknieModalnym } from './pytanie-modalne.ts';

const NAGLOWEK_MOBILE = 'Mobile';
const NAGLOWEK_AOD = 'Always On Display';

const STEROWANIA: ReadonlyArray<readonly [string, string]> = [
  [MobileProcessControl.Pause, 'wstrzymaj'],
  [MobileProcessControl.Resume, 'wznów'],
  [MobileProcessControl.Stop, 'zatrzymaj'],
  [MobileProcessControl.Restart, 'uruchom ponownie'],
  [MobileProcessControl.Approve, 'zatwierdź krok'],
];

export function zwiazNawiazania(kanal: Kanal, korzen: Element): void {
  const sekcja = korzen.querySelector('#sekcja-nawiazania');
  if (sekcja === null) return;
  const karty = sekcja.querySelectorAll('.us-nawiazanie');
  const mobile = karty[0];
  const aod = karty[1];
  if (mobile !== undefined) dolozCzynnosciMobile(mobile);
  if (aod !== undefined) dolozCzynnosciAod(aod);

  sekcja.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const czynnosc = cel.closest<HTMLElement>('[data-nawiazanie]');
    if (czynnosc === null) {
      if (cel.closest('#us-aod-przelacz') !== null) {
        zdarzenie.stopPropagation();
        void przelaczNakladke(kanal, sekcja);
      }
      return;
    }
    zdarzenie.stopPropagation();
    void wykonajCzynnosc(kanal, sekcja, czynnosc.dataset.nawiazanie ?? '');
  });

  void odczytajStanMobile(kanal, sekcja);
  void odczytajStanAod(kanal, sekcja);
}

/* Karty Nawiązań są w makiecie odesłaniami; czynności dokłada wiązanie, bo bez
   nich Mobile i nakładka nie mają w Ustawieniach żadnej drogi do rdzenia. */
function dolozCzynnosciMobile(karta: Element): void {
  const wiersz = karta.querySelector('.us-wiersz');
  if (wiersz === null) return;
  wiersz.append(
    przycisk(karta.ownerDocument, 'mobile-stan', 'Odczytaj stan'),
    przycisk(karta.ownerDocument, 'mobile-procesy', 'Procesy telefonu'),
    przycisk(karta.ownerDocument, 'mobile-steruj', 'Steruj procesem'),
  );
}

function dolozCzynnosciAod(karta: Element): void {
  const wiersz = karta.querySelector('.us-wiersz');
  if (wiersz === null) return;
  const dokument = karta.ownerDocument;
  wiersz.append(
    przycisk(dokument, 'aod-kontekst', 'Kontekst okna'),
    przycisk(dokument, 'aod-podpowiedzi', 'Podpowiedzi'),
    przycisk(dokument, 'aod-wiadomosc', 'Wyślij wiadomość'),
    przycisk(dokument, 'aod-glos', 'Polecenie głosowe'),
    przycisk(dokument, 'aod-przypnij', 'Przypnij proces'),
    przycisk(dokument, 'aod-odepnij', 'Odepnij proces'),
    przycisk(dokument, 'aod-sygnaly', 'Sygnały'),
    przycisk(dokument, 'aod-zglos', 'Zgłoś sygnał'),
  );
}

function przycisk(dokument: Document, kod: string, napis: string): HTMLButtonElement {
  const guzik = dokument.createElement('button');
  guzik.type = 'button';
  guzik.className = 'dn-btn dn-btn--duch dn-btn--sm';
  guzik.dataset.nawiazanie = kod;
  guzik.textContent = napis;
  return guzik;
}

async function wykonajCzynnosc(
  kanal: Kanal,
  sekcja: Element,
  kod: string,
): Promise<void> {
  if (kod === 'mobile-stan') return odczytajStanMobile(kanal, sekcja);
  if (kod === 'mobile-procesy') return wypiszProcesy(kanal);
  if (kod === 'mobile-steruj') return sterujProcesem(kanal);
  if (kod === 'aod-kontekst') return pobierzKontekst(kanal);
  if (kod === 'aod-podpowiedzi') return pobierzPodpowiedzi(kanal);
  if (kod === 'aod-wiadomosc') return wyslijWiadomosc(kanal);
  if (kod === 'aod-glos') return wydajPolecenieGlosem(kanal);
  if (kod === 'aod-przypnij') return przypnijProces(kanal, sekcja, true);
  if (kod === 'aod-odepnij') return przypnijProces(kanal, sekcja, false);
  if (kod === 'aod-sygnaly') return wypiszSygnaly(kanal);
  if (kod === 'aod-zglos') return zglosSygnal(kanal);
  oglos(NAGLOWEK_MOBILE, 'Ta czynność nie ma wiązania.', 'ostrzezenie');
}

/* Stan mobilny liczy się bez wskazania urządzenia: Ustawienia stoją na
   pulpicie, a parowanie telefonu jest osobną sprawą sekcji Urządzenia. */
async function odczytajStanMobile(kanal: Kanal, sekcja: Element): Promise<void> {
  const wynik = await wywolaj(kanal, Command.MobileStatusGet, {});
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK_MOBILE, wynik.blad?.message ?? 'Stan mobilny nie doszedł.',
      'ostrzezenie');
    return;
  }
  const stan = wynik.wynik.status;
  const opis = sekcja.querySelector('.us-nawiazanie .us-opis');
  const zdanie = `sesji: ${String(stan.sessionCount)}`
    + ` · okien: ${String(stan.windowCount)}`
    + ` · procesów w biegu: ${String(stan.runningProcessCount)}`
    + ` · ${stan.paired ? 'urządzenie sparowane' : 'bez sparowanego urządzenia'}`;
  if (opis !== null) opis.textContent = zdanie;
  oglos(NAGLOWEK_MOBILE, zdanie);
}

async function wypiszProcesy(kanal: Kanal): Promise<void> {
  const wynik = await wywolaj(kanal, Command.MobileProcessList, {});
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK_MOBILE, wynik.blad?.message ?? 'Wykaz procesów nie doszedł.',
      'ostrzezenie');
    return;
  }
  const procesy = wynik.wynik.processes;
  if (procesy.length === 0) {
    oglos(NAGLOWEK_MOBILE, 'Rdzeń nie prowadzi teraz żadnego procesu.');
    return;
  }
  oglos(NAGLOWEK_MOBILE, `Procesów: ${String(procesy.length)}; `
    + procesy.slice(0, 3).map((poz) => `${poz.label} (${poz.status})`).join(', '));
}

/* Sterowanie idzie po identyfikatorze procesu, więc pytanie podaje wykaz
   biegnących: wpisanie identyfikatora z pamięci nie jest wyborem. */
async function sterujProcesem(kanal: Kanal): Promise<void> {
  const spis = await wywolaj(kanal, Command.MobileProcessList, {});
  if (!spis.udany || spis.wynik === undefined) {
    oglos(NAGLOWEK_MOBILE, spis.blad?.message ?? 'Wykaz procesów nie doszedł.',
      'ostrzezenie');
    return;
  }
  const procesy = spis.wynik.processes;
  if (procesy.length === 0) {
    oglos(NAGLOWEK_MOBILE, 'Nie ma czym sterować: wykaz procesów jest pusty.');
    return;
  }
  const odpowiedz = await zapytajWOknieModalnym({
    tytul: 'Sterowanie procesem z telefonu',
    opis: 'Zatrzymanie przerywa bieg; pozostałe sterowania są odwracalne.',
    pola: [
      {
        klucz: 'proces',
        etykieta: 'Proces',
        wybor: procesy.map((poz) => [poz.id, `${poz.label} (${poz.status})`] as const),
        wymagane: true,
      },
      {
        klucz: 'sterowanie',
        etykieta: 'Sterowanie',
        wybor: STEROWANIA,
        wymagane: true,
      },
    ],
    wykonanie: 'Wydaj sterowanie',
  });
  if (odpowiedz === null) {
    oglos(NAGLOWEK_MOBILE, 'Sterowanie odwołane; procesy bez zmiany.');
    return;
  }
  const wynik = await wywolaj(kanal, Command.MobileProcessControl, {
    processId: odpowiedz['proces'] ?? '',
    control: (odpowiedz['sterowanie'] ?? '') as typeof MobileProcessControl.Pause,
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK_MOBILE, wynik.blad?.message ?? 'Rdzeń odmówił sterowania.',
      'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK_MOBILE, `Sterowanie przyjęte; stan procesu: `
    + `${wynik.wynik?.process.status ?? 'bez wskazania'}.`);
}

async function odczytajStanAod(kanal: Kanal, sekcja: Element): Promise<void> {
  const wynik = await wywolaj(kanal, Command.AodStatusGet, {});
  const plakietka = sekcja.querySelector('[data-stan-aod]');
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK_AOD, wynik.blad?.message ?? 'Stan nakładki nie doszedł.',
      'ostrzezenie');
    return;
  }
  const stan = wynik.wynik.status;
  const przypiete = stan.attachedProcessIds?.length ?? 0;
  if (plakietka !== null) {
    plakietka.textContent = `w biegu: ${String(stan.runningProcessCount)}`
      + ` · przypiętych: ${String(przypiete)}`;
  }
  oglos(NAGLOWEK_AOD, `Nakładka zna ${String(stan.runningProcessCount)} procesów`
    + ` w biegu, przypiętych ${String(przypiete)}.`);
}

/* Przełącznik makiety nie miał czego przełączać: wyciszenie czasowe jest
   jedyną nastawą nakładki, którą Ustawienia zakładają i znoszą jednym ruchem. */
async function przelaczNakladke(kanal: Kanal, sekcja: Element): Promise<void> {
  const stojace = await wywolaj(kanal, Command.AodMuteGet, {});
  if (!stojace.udany || stojace.wynik === undefined) {
    oglos(NAGLOWEK_AOD, stojace.blad?.message ?? 'Wykaz wyciszeń nie doszedł.',
      'ostrzezenie');
    return;
  }
  const czasowe = stojace.wynik.mutes.find((poz) => poz.kind === AodMuteKind.Timed);
  if (czasowe !== undefined) {
    await zniesWyciszenie(kanal, sekcja, czasowe.id);
    return;
  }
  await zalozWyciszenie(kanal, sekcja);
}

async function zniesWyciszenie(
  kanal: Kanal,
  sekcja: Element,
  idWyciszenia: string,
): Promise<void> {
  const wynik = await wywolaj(kanal, Command.AodMuteSet, {
    muteId: idWyciszenia,
    muted: false,
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK_AOD, wynik.blad?.message ?? 'Rdzeń odmówił zniesienia ciszy.',
      'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK_AOD, 'Wyciszenie zniesione; nakładka znów się odzywa.');
  await odczytajStanAod(kanal, sekcja);
}

/* Rdzeń odrzuca wyciszenie czasowe bez chwili końca, i słusznie: cisza bez
   terminu nie mówi Operatorowi, do kiedy trwa. Więc pytamy o czas trwania. */
async function zalozWyciszenie(kanal: Kanal, sekcja: Element): Promise<void> {
  const odpowiedz = await zapytajWOknieModalnym({
    tytul: 'Wyciszenie nakładki na czas',
    opis: 'Po upływie czasu nakładka odzywa się sama; można ją odciszyć wcześniej.',
    pola: [{
      klucz: 'minuty',
      etykieta: 'Czas ciszy',
      wybor: [
        ['15', 'kwadrans'],
        ['60', 'godzina'],
        ['480', 'osiem godzin'],
      ],
      wymagane: true,
    }],
    wykonanie: 'Wycisz',
  });
  if (odpowiedz === null) {
    oglos(NAGLOWEK_AOD, 'Wyciszenie odwołane; nakładka odzywa się dalej.');
    return;
  }
  const minuty = Number(odpowiedz['minuty'] ?? '');
  if (!Number.isFinite(minuty) || minuty <= 0) {
    oglos(NAGLOWEK_AOD, 'Czas ciszy musi być liczbą minut.', 'ostrzezenie');
    return;
  }
  const koniec = Date.now() + minuty * 60_000;
  const wynik = await wywolaj(kanal, Command.AodMuteSet, {
    kind: AodMuteKind.Timed,
    endsAt: koniec,
    muted: true,
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK_AOD, wynik.blad?.message ?? 'Rdzeń odmówił wyciszenia.',
      'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK_AOD, `Nakładka wyciszona na ${String(minuty)} minut;`
    + ` wyciszeń czynnych: ${String(wynik.wynik?.mutes.length ?? 0)}.`);
  await odczytajStanAod(kanal, sekcja);
}

async function pobierzKontekst(kanal: Kanal): Promise<void> {
  const wynik = await wywolaj(kanal, Command.AodContextGet, { historyLimit: 10 });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK_AOD, wynik.blad?.message ?? 'Kontekst nie doszedł.', 'ostrzezenie');
    return;
  }
  const komplet = wynik.wynik.context;
  oglos(NAGLOWEK_AOD, 'Kontekst okna: historia '
    + `${String(komplet.historyMessageIds?.length ?? 0)} wiadomości`
    + `, dokumentów ${String(komplet.documentIds?.length ?? 0)}`
    + `, agentów ${String(komplet.agentIds?.length ?? 0)}.`);
}

async function pobierzPodpowiedzi(kanal: Kanal): Promise<void> {
  const wynik = await wywolaj(kanal, Command.AodSuggestion, { limit: 5 });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK_AOD, wynik.blad?.message ?? 'Podpowiedzi nie doszły.', 'ostrzezenie');
    return;
  }
  const podpowiedzi = wynik.wynik.suggestions;
  oglos(NAGLOWEK_AOD, podpowiedzi.length === 0
    ? 'Nakładka nie ma teraz nic do podpowiedzenia.'
    : `Podpowiedzi: ${String(podpowiedzi.length)}; pierwsza: `
      + `${podpowiedzi[0]?.text ?? 'bez treści'}`);
}

async function wyslijWiadomosc(kanal: Kanal): Promise<void> {
  const odpowiedz = await zapytajWOknieModalnym({
    tytul: 'Wiadomość z nakładki',
    opis: 'Bez wskazania okna wiadomość idzie do okna ogniskowanego.',
    pola: [{ klucz: 'tresc', etykieta: 'Treść', obszerne: true, wymagane: true }],
    wykonanie: 'Wyślij',
  });
  if (odpowiedz === null) {
    oglos(NAGLOWEK_AOD, 'Wysyłka odwołana; nic nie poszło do okna.');
    return;
  }
  const wynik = await wywolaj(kanal, Command.AodChatSend, {
    text: odpowiedz['tresc'] ?? '',
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK_AOD, wynik.blad?.message ?? 'Rdzeń odmówił wysyłki.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK_AOD, 'Wiadomość przyjęta przez okno rozmowy.');
}

/* Nakładka przyjmuje głos albo poprawioną transkrypcję; Ustawienia nie mają
   mikrofonu, więc podają samą transkrypcję. */
async function wydajPolecenieGlosem(kanal: Kanal): Promise<void> {
  const odpowiedz = await zapytajWOknieModalnym({
    tytul: 'Polecenie głosowe nakładki',
    opis: 'Treść wchodzi jako transkrypcja poprawiona przez Operatora.',
    pola: [{ klucz: 'tresc', etykieta: 'Transkrypcja', wymagane: true }],
    wykonanie: 'Wydaj polecenie',
  });
  if (odpowiedz === null) {
    oglos(NAGLOWEK_AOD, 'Polecenie odwołane; asystent nic nie dostał.');
    return;
  }
  const wynik = await wywolaj(kanal, Command.AodVoiceCommand, {
    transcript: odpowiedz['tresc'] ?? '',
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK_AOD, wynik.blad?.message ?? 'Rdzeń odmówił polecenia.',
      'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK_AOD, 'Polecenie przyjęte przez asystenta.');
}

async function przypnijProces(
  kanal: Kanal,
  sekcja: Element,
  przypiecie: boolean,
): Promise<void> {
  const spis = await wywolaj(kanal, Command.MobileProcessList, {});
  const procesy = spis.udany ? (spis.wynik?.processes ?? []) : [];
  if (procesy.length === 0) {
    oglos(NAGLOWEK_AOD, 'Wykaz procesów jest pusty; nie ma czego przypinać.');
    return;
  }
  const odpowiedz = await zapytajWOknieModalnym({
    tytul: przypiecie ? 'Przypnij proces do nakładki' : 'Odepnij proces od nakładki',
    pola: [{
      klucz: 'proces',
      etykieta: 'Proces',
      wybor: procesy.map((poz) => [poz.id, poz.label] as const),
      wymagane: true,
    }],
    wykonanie: przypiecie ? 'Przypnij' : 'Odepnij',
  });
  if (odpowiedz === null) {
    oglos(NAGLOWEK_AOD, 'Czynność odwołana; obserwacja bez zmiany.');
    return;
  }
  const proces = odpowiedz['proces'] ?? '';
  const wynik = przypiecie
    ? await wywolaj(kanal, Command.AodObserveAttach, { processId: proces })
    : await wywolaj(kanal, Command.AodObserveDetach, { processId: proces });
  if (!wynik.udany) {
    oglos(NAGLOWEK_AOD, wynik.blad?.message ?? 'Rdzeń odmówił zmiany obserwacji.',
      'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK_AOD, `Przypiętych procesów: `
    + `${String(wynik.wynik?.attachedProcessIds.length ?? 0)}.`);
  await odczytajStanAod(kanal, sekcja);
}

async function wypiszSygnaly(kanal: Kanal): Promise<void> {
  const wynik = await wywolaj(kanal, Command.AodSignalList, { limit: 10 });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK_AOD, wynik.blad?.message ?? 'Sygnały nie doszły.', 'ostrzezenie');
    return;
  }
  const sygnaly = wynik.wynik.signals;
  oglos(NAGLOWEK_AOD, sygnaly.length === 0
    ? 'Nakładka nie ma zapisanych sygnałów.'
    : `Sygnałów: ${String(sygnaly.length)}; ostatni: `
      + `${sygnaly[0]?.text ?? 'bez treści'}`);
}

async function zglosSygnal(kanal: Kanal): Promise<void> {
  const odpowiedz = await zapytajWOknieModalnym({
    tytul: 'Zgłoszenie sygnału do nakładki',
    opis: 'Klasa zdarzenia rozstrzyga, czy sygnał wpadnie w wyciszenie.',
    pola: [
      {
        klucz: 'klasa',
        etykieta: 'Klasa zdarzenia',
        wybor: [
          [AodEventClass.ExecutionLoopState, 'stan pętli wykonawczej'],
          [AodEventClass.TaskQueueState, 'stan kolejki zadań'],
          [AodEventClass.QualityControlResult, 'wynik kontroli jakości'],
          [AodEventClass.ModuleEvent, 'zdarzenie modułu'],
          [AodEventClass.Schedule, 'harmonogram'],
        ],
        wymagane: true,
      },
      { klucz: 'tresc', etykieta: 'Treść sygnału', wymagane: true },
    ],
    wykonanie: 'Zgłoś sygnał',
  });
  if (odpowiedz === null) {
    oglos(NAGLOWEK_AOD, 'Zgłoszenie odwołane; nakładka nic nie dostała.');
    return;
  }
  const wynik = await wywolaj(kanal, Command.AodSignalReport, {
    eventClass: (odpowiedz['klasa'] ?? '') as typeof AodEventClass.ModuleEvent,
    text: odpowiedz['tresc'] ?? '',
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK_AOD, wynik.blad?.message ?? 'Rdzeń odmówił zgłoszenia.',
      'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK_AOD, 'Sygnał przyjęty przez nakładkę.');
}
