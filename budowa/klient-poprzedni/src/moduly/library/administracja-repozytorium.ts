import type { LibraryRule, LibraryWebhook } from '../../../../shared/contract';
import { opisOdmowy } from '../../komponenty/odmowa';
import { poleTekstowe, przycisk, utworzWierszOdpowiedzi } from '../../modele/kontrolki-formularza';
import type { StanBiblioteki } from './stan-biblioteki';

/**
 * Warstwa czwarta modułu Library: zdolności eksperckie repozytorium niewidoczne dla użytkownika
 * podstawowego — reguły, schemat metadanych, retencja, udostępnienia, nasłuchy, klasyfikacja
 * i porównanie zasobów.
 */
export interface AdministracjaRepozytorium {
  element: HTMLElement;
  /** Odczytuje wykazy warstwy czwartej: reguły, nasłuchy, udostępnienia, sugestie. */
  odswiez(): void;
}

export function utworzAdministracjeRepozytorium(stan: StanBiblioteki): AdministracjaRepozytorium {
  const odpowiedz = utworzWierszOdpowiedzi();

  const wykazy = document.createElement('div');
  wykazy.className = 'ml-administracja__wykazy';

  /** Wypisuje wykaz jako listę pozycji wraz z nagłówkiem. */
  function blokWykazu(tytul: string, wiersze: readonly string[]): HTMLElement {
    const naglowek = document.createElement('strong');
    naglowek.className = 'ml-higiena__tytul';
    naglowek.textContent = `${tytul}: ${wiersze.length}`;

    const lista = document.createElement('ul');
    lista.className = 'ml-higiena__lista';
    for (const wiersz of wiersze) {
      const pozycja = document.createElement('li');
      pozycja.textContent = wiersz;
      lista.append(pozycja);
    }
    const blok = document.createElement('div');
    blok.className = 'ml-administracja__blok';
    blok.append(naglowek, lista);
    return blok;
  }

  // ── Reguła kolekcji inteligentnej ─────────────────────────────────────────
  const nazwaReguly = poleTekstowe({
    etykieta: 'Nazwa reguły',
    podpowiedz: 'np. materiały rynku X',
  });
  const etykietaWarunku = poleTekstowe({
    etykieta: 'Warunek: etykieta',
    podpowiedz: 'np. rynek-x',
    opis: 'Zasób noszący tę etykietę trafia do kolekcji docelowej przy przeliczeniu reguły.',
  });
  const kolekcjaReguly = poleTekstowe({
    etykieta: 'Kolekcja docelowa (identyfikator)',
    podpowiedz: 'kolek-…',
  });

  const zapiszRegule = przycisk('Zapisz regułę i przelicz', 'dn-btn dn-btn--sm dn-btn--zarys');
  zapiszRegule.dataset['czynnosc'] = 'regula-zapis';
  zapiszRegule.addEventListener('click', () => void ustawRegule());

  async function ustawRegule(): Promise<void> {
    const regula = {
      id: '',
      kind: 'collection',
      name: nazwaReguly.kontrolka.value.trim(),
      condition: { tags: [etykietaWarunku.kontrolka.value.trim()] },
      targetCollectionId: kolekcjaReguly.kontrolka.value.trim(),
      enabled: true,
      createdAt: 0,
    } as unknown as LibraryRule;

    odpowiedz.pokaz('Zapisywanie reguły i przeliczanie kolekcji…', true);
    const wynik = await stan.zrodlo.zapiszRegule(regula, true);
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(opisOdmowy('Reguła', wynik.blad?.code, wynik.blad?.message), false);
      return;
    }
    odpowiedz.pokaz(
      `Reguła ${wynik.wynik.rule.name} zapisana. Zasobów spełniających warunek: ` +
        `${wynik.wynik.matchedFiles ?? 0}.`,
      true,
    );
    await odczytajWykazy();
    await stan.odczytaj('', '');
  }

  // ── Pole niestandardowe schematu metadanych ───────────────────────────────
  const kodPola = poleTekstowe({ etykieta: 'Kod pola schematu', podpowiedz: 'np. numerSprawy' });
  const nazwaPola = poleTekstowe({ etykieta: 'Nazwa pola widoczna dla Operatora' });

  const zapiszPole = przycisk('Załóż pole schematu', 'dn-btn dn-btn--sm dn-btn--zarys');
  zapiszPole.dataset['czynnosc'] = 'schemat-zapis';
  zapiszPole.addEventListener('click', () => void ustawPoleSchematu());

  async function ustawPoleSchematu(): Promise<void> {
    odpowiedz.pokaz('Zapisywanie definicji pola schematu…', true);
    const wynik = await stan.zrodlo.zapiszPoleSchematu(
      {
        code: kodPola.kontrolka.value.trim(),
        label: nazwaPola.kontrolka.value.trim(),
        kind: 'text',
        required: false,
        createdAt: 0,
      },
      false,
    );
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(opisOdmowy('Schemat metadanych', wynik.blad?.code, wynik.blad?.message), false);
      return;
    }
    const schemat = await stan.zrodlo.schematMetadanych();
    const ile = schemat.udany && schemat.wynik !== undefined ? schemat.wynik.fields.length : 0;
    odpowiedz.pokaz(
      `Pole ${wynik.wynik.field.code} zapisane. Zasobów z wartością tego pola: ` +
        `${wynik.wynik.affectedFiles}. Pól w schemacie: ${ile}.`,
      true,
    );
  }

  // ── Polityka retencji ─────────────────────────────────────────────────────
  const dniRetencji = poleTekstowe({
    etykieta: 'Okres przechowywania (dni)',
    podpowiedz: 'np. 365',
    opis: 'Po upływie okresu zasób zostaje ZGŁOSZONY do przeglądu; nic nie znika samo.',
  });

  const zapiszRetencje = przycisk('Ustaw politykę retencji', 'dn-btn dn-btn--sm dn-btn--zarys');
  zapiszRetencje.dataset['czynnosc'] = 'retencja-zapis';
  zapiszRetencje.addEventListener('click', () => void ustawRetencje());

  async function ustawRetencje(): Promise<void> {
    const dni = Number.parseInt(dniRetencji.kontrolka.value, 10);
    if (!Number.isFinite(dni) || dni <= 0) {
      odpowiedz.pokaz('Okres przechowywania podaje się w dniach, liczbą dodatnią.', false);
      return;
    }
    odpowiedz.pokaz('Zapisywanie polityki przechowywania…', true);
    const wynik = await stan.zrodlo.zapiszRetencje(
      { id: '', scope: 'global', keepDays: dni, action: 'review', createdAt: 0 },
      false,
    );
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(opisOdmowy('Retencja', wynik.blad?.code, wynik.blad?.message), false);
      return;
    }
    odpowiedz.pokaz(
      `Polityka zapisana: ${wynik.wynik.policy.keepDays} dni, czynność ` +
        `${wynik.wynik.policy.action}. Zasobów objętych: ${wynik.wynik.affectedFiles}.`,
      true,
    );
    await odczytajWykazy();
  }

  // ── Nasłuch zewnętrzny ────────────────────────────────────────────────────
  const adresNasluchu = poleTekstowe({
    etykieta: 'Adres nasłuchu (HTTP/HTTPS)',
    podpowiedz: 'https://…',
    opis: 'Rdzeń zgłasza pod ten adres zdarzenia repozytorium; sekret podpisu jest jawny.',
  });

  const zapiszNasluch = przycisk('Załóż nasłuch', 'dn-btn dn-btn--sm dn-btn--zarys');
  zapiszNasluch.dataset['czynnosc'] = 'nasluch-zapis';
  zapiszNasluch.addEventListener('click', () => void ustawNasluch());

  async function ustawNasluch(): Promise<void> {
    const nasluch = {
      id: '',
      url: adresNasluchu.kontrolka.value.trim(),
      events: ['fileAdded', 'fileChanged', 'fileArchived'],
      enabled: true,
      createdAt: 0,
    } as unknown as LibraryWebhook;

    odpowiedz.pokaz('Zapisywanie nasłuchu zewnętrznego…', true);
    const wynik = await stan.zrodlo.zapiszNasluch(nasluch);
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(opisOdmowy('Nasłuch', wynik.blad?.code, wynik.blad?.message), false);
      return;
    }
    odpowiedz.pokaz(
      `Nasłuch ${wynik.wynik.webhook.id} zapisany dla zdarzeń: ` +
        `${wynik.wynik.webhook.events.join(', ')}.`,
      true,
    );
    await odczytajWykazy();
  }

  // ── Klasyfikacja wsadowa i sugestie ───────────────────────────────────────
  const klasyfikuj = przycisk('Klasyfikuj zaznaczenie', 'dn-btn dn-btn--sm dn-btn--zarys');
  klasyfikuj.dataset['czynnosc'] = 'klasyfikacja';
  klasyfikuj.addEventListener('click', () => void klasyfikujZbior());

  async function klasyfikujZbior(): Promise<void> {
    odpowiedz.pokaz('Rdzeń klasyfikuje zbiór i składa sugestie porządkujące…', true);
    const wynik = await stan.zrodlo.klasyfikuj(stan.zaznaczone(), false);
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(opisOdmowy('Klasyfikacja', wynik.blad?.code, wynik.blad?.message), false);
      return;
    }
    odpowiedz.pokaz(
      `Przetworzono zasobów: ${wynik.wynik.processedCount}, sugestii: ` +
        `${wynik.wynik.suggestions.length}. Sugestie czekają na decyzję — zapisu bez ` +
        'pytania moduł nie robi.',
      true,
    );
    await odczytajWykazy();
  }

  const przyjmijSugestie = przycisk('Przyjmij wszystkie sugestie', 'dn-btn dn-btn--sm dn-btn--zarys');
  przyjmijSugestie.dataset['czynnosc'] = 'sugestie-przyjecie';
  przyjmijSugestie.addEventListener('click', () => void rozstrzygnij(true));

  const odrzucSugestie = przycisk('Odrzuć wszystkie sugestie', 'dn-btn dn-btn--sm dn-btn--zarys');
  odrzucSugestie.dataset['czynnosc'] = 'sugestie-odrzucenie';
  odrzucSugestie.addEventListener('click', () => void rozstrzygnij(false));

  async function rozstrzygnij(przyjmij: boolean): Promise<void> {
    const wykaz = await stan.zrodlo.wykazSugestii();
    if (!wykaz.udany || wykaz.wynik === undefined || wykaz.wynik.suggestions.length === 0) {
      odpowiedz.pokaz('Ani jednej sugestii oczekującej na decyzję.', false);
      return;
    }
    const kody = wykaz.wynik.suggestions.map((sugestia) => sugestia.id);
    odpowiedz.pokaz(przyjmij ? 'Przyjmowanie sugestii…' : 'Odrzucanie sugestii…', true);
    const wynik = await stan.zrodlo.rozstrzygnijSugestie(kody, przyjmij);
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(opisOdmowy('Sugestie', wynik.blad?.code, wynik.blad?.message), false);
      return;
    }
    odpowiedz.pokaz(
      `Przyjętych: ${wynik.wynik.appliedCount}, odrzuconych: ${wynik.wynik.rejectedCount}. ` +
        'Przyjęcie WYKONUJE czynność, którą sugestia opisuje — nie samo ją odhacza.',
      true,
    );
    await odczytajWykazy();
    await stan.odczytaj('', '');
  }

  // ── Porównanie dwóch zasobów ──────────────────────────────────────────────
  const porownaj = przycisk('Porównaj dwa zaznaczone', 'dn-btn dn-btn--sm dn-btn--zarys');
  porownaj.dataset['czynnosc'] = 'porownanie';
  porownaj.addEventListener('click', () => void porownajZaznaczone());

  async function porownajZaznaczone(): Promise<void> {
    const zaznaczone = stan.zaznaczone();
    if (zaznaczone.length !== 2) {
      odpowiedz.pokaz('Porównanie obejmuje dwa zasoby — zaznacz dokładnie dwa pliki.', false);
      return;
    }
    odpowiedz.pokaz('Rdzeń porównuje treści zasobów…', true);
    const wynik = await stan.zrodlo.porownaj(zaznaczone[0] ?? '', zaznaczone[1] ?? '');
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(opisOdmowy('Porównanie', wynik.blad?.code, wynik.blad?.message), false);
      return;
    }
    const roznica = wynik.wynik.diff;
    wykazy.replaceChildren(
      blokWykazu(
        `Różnice: ${roznica.leftLabel} ↔ ${roznica.rightLabel}`,
        roznica.hunks.map(
          (fragment) =>
            `${fragment.kind} · wiersze ${fragment.leftFrom}-${fragment.leftTo} → ` +
            `${fragment.rightFrom}-${fragment.rightTo}: ${fragment.text}`,
        ),
      ),
    );
    odpowiedz.pokaz(
      roznica.identical
        ? 'Treści są identyczne.'
        : `Różnic: ${roznica.hunks.length}.` +
            (roznica.comparedAsText
              ? ' Porównano tekst wydobyty z dokumentu, nie bajty — dla dokumentu binarnego ' +
                'jest to jedyna droga.'
              : ''),
      true,
    );
  }

  /** Odczytuje wykazy warstwy czwartej z rdzenia. */
  async function odczytajWykazy(): Promise<void> {
    const [reguly, nasluchy, udostepnienia, sugestie] = await Promise.all([
      stan.zrodlo.wykazRegul(),
      stan.zrodlo.wykazNasluchow(),
      stan.zrodlo.wykazUdostepnien(),
      stan.zrodlo.wykazSugestii(),
    ]);
    const bloki: HTMLElement[] = [];
    if (reguly.udany && reguly.wynik !== undefined) {
      bloki.push(
        blokWykazu(
          'Reguły repozytorium',
          reguly.wynik.rules.map(
            (regula) =>
              `${regula.name} (${regula.kind}) → ${regula.targetCollectionId ?? 'bez kolekcji'}` +
              `${regula.enabled ? '' : ' — wyłączona'}`,
          ),
        ),
      );
    }
    if (nasluchy.udany && nasluchy.wynik !== undefined) {
      bloki.push(
        blokWykazu(
          'Nasłuchy zewnętrzne',
          nasluchy.wynik.webhooks.map(
            (nasluch) => `${nasluch.url} — ${nasluch.events.join(', ')}`,
          ),
        ),
      );
    }
    if (udostepnienia.udany && udostepnienia.wynik !== undefined) {
      bloki.push(
        blokWykazu(
          'Udostępnienia',
          udostepnienia.wynik.shares.map(
            (udostepnienie) =>
              `${udostepnienie.scope}: ${udostepnienie.targetId} — token ` +
              `${udostepnienie.token.slice(0, 12)}…` +
              (udostepnienie.revokedAt === undefined ? '' : ' (odwołane)'),
          ),
        ),
      );
    }
    if (sugestie.udany && sugestie.wynik !== undefined) {
      bloki.push(
        blokWykazu(
          'Sugestie oczekujące',
          sugestie.wynik.suggestions.map(
            (sugestia) =>
              `${sugestia.kind} · ${sugestia.fileId}` +
              `${sugestia.value === undefined ? '' : ` → ${sugestia.value}`} — ${sugestia.reason}`,
          ),
        ),
      );
    }
    wykazy.replaceChildren(...bloki);
  }

  const pasek = document.createElement('div');
  pasek.className = 'ml-administracja__pasek';
  pasek.append(
    zapiszRegule,
    zapiszPole,
    zapiszRetencje,
    zapiszNasluch,
    klasyfikuj,
    przyjmijSugestie,
    odrzucSugestie,
    porownaj,
  );

  const naglowek = document.createElement('h4');
  naglowek.className = 'ml-metadane__naglowek';
  naglowek.textContent = 'Sterowanie repozytorium — warstwa ekspercka';

  const element = document.createElement('div');
  element.className = 'ml-administracja';
  element.append(
    naglowek,
    nazwaReguly.element,
    etykietaWarunku.element,
    kolekcjaReguly.element,
    kodPola.element,
    nazwaPola.element,
    dniRetencji.element,
    adresNasluchu.element,
    pasek,
    wykazy,
    odpowiedz.element,
  );

  return {
    element,
    odswiez() {
      void odczytajWykazy();
    },
  };
}
