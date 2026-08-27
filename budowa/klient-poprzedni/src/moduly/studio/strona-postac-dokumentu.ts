import './strona-postac.css';

import {
  StudioAuthor,
  StudioTabKind,
  StudioTabLeader,
  type StudioActionBalance,
  type StudioCharacterFormat,
  type StudioDocumentForm,
  type StudioParagraphFormat,
  type StudioProvenance,
} from '../../../../shared/contract';
import { opisOdmowy } from '../../komponenty/odmowa';
import { przycisk, utworzWierszOdpowiedzi } from '../../modele/kontrolki-formularza';
import type { RodzajTabulatora, TabulatorAkapitu, ZnakWiodacy } from './linijka-podzialka';
import { wchlonNosnikiKontraktu, type StronaPracy } from './nastawy-strony';
import type { Zaznaczenie } from './pola-stanu';
import { utworzMagazynWidokuRdzenia, type MagazynWidokuRdzenia } from './strona-magazyn-widoku';
import { utworzStronePanelNastaw, type StronaPanelNastaw } from './strona-panel-nastaw';
import { utworzStronePanelTresci, type StronaPanelTresci } from './strona-panel-tresci';
import { opiszBilans } from './strona-pola-postaci';
import { utworzStylPanelArkusza, type StylPanelArkusza } from './styl-panel-arkusza';
import { utworzStylPanelList, type StylPanelList } from './styl-panel-list';
import type { ZrodloPostaciStudio } from './zrodlo-postaci-studio';


/** Zaplecze warstwy postaci: źródło rdzenia, zaznaczenie bieżące, autor czynności i sposób ogłoszenia zdania — wszystko, czego panel sam nie zna. */
export interface ZapleczePostaci {
  zrodlo: ZrodloPostaciStudio;
  /** Dokument czynny; pusty znaczy „żadnego nie wczytano". */
  idDokumentu(): string;
  /** Okno komunikacji sesji; puste przed osadzeniem. */
  idOkna(): string;
  /** Zaznaczenie w treści; `null` znaczy zakres „cały dokument". */
  zaznaczenie(): Zaznaczenie | null;
  /** Miejsce kursora w znakach — podział i znak wstawiają się tam. */
  miejsceKursora(): number;
  /** Postać zmieniona przez rdzeń — okno przerysowuje nią kartkę bez ponownego wczytania dokumentu. */
  naPostac(postac: StudioDocumentForm): void;
  /** Zdanie o skutku czynności — pasek odpowiedzi okna. */
  naZdanie(tresc: string, powodzenie: boolean): void;
}

/** Warstwa postaci dokumentu wraz z panelami i drogami zgłaszania zmian z linijki do rdzenia oraz zapisu i odczytu postaci dla dokumentu. */
export interface PostacDokumentu {
  /** Pas przycisków otwierających trzy panele wraz z panelami. */
  element: HTMLElement;
  /** Odczyt wszystkiego, co panele pokazują: nastawy, sekcje, style, listy, znaki. */
  wczytaj(): Promise<void>;
  /** Powtórny odczyt po zmianie dokumentu czynnego. */
  odswiez(): void;

  /* ── Drogi dla linijek i wstążki ─────────────────────────────────────────── */

  /** Margines przestawiony chwytem na linijce jedzie tą samą komendą do rdzenia co pole panelu. */
  zglosMargines(ktory: 'gora' | 'dol' | 'lewy' | 'prawy', milimetry: number): void;
  /** Powierzchnia oddaje całą StronaPracy jedną komendą, nie nazwę zmienionego pola. */
  zglosStrone(strona: StronaPracy): void;
  /** Wcięcia akapitu przestawione chwytem — idą stylem akapitu fragmentu. */
  zglosWciecia(wciecia: {
    pierwszyWierszMm?: number;
    leweMm?: number;
    praweMm?: number;
  }): void;
  /** Tabulator założony, przestawiony albo zdjęty chwytem na linijce. */
  zglosTabulator(polozenieMm: number, zdejmij: boolean): void;
  /** Linijka oddaje cały wykaz tabulatorów, warstwa wysyła rdzeniowi tylko różnicę. */
  zglosTabulatory(tabulatory: readonly TabulatorAkapitu[]): void;

  /* ── Trwałość postaci i treści ───────────────────────────────────────────── */

  /** Odczyt pełnej postaci dokumentu — arkusz stylów, sekcje, bloki, pola. */
  odczytajPostac(): Promise<StudioDocumentForm | null>;
  /** Ta komenda, inaczej niż zapis napisu i tytułu, przenosi arkusz stylów, strony, sekcje i tabele. */
  zapiszPostac(postac: StudioDocumentForm, tresc: string, tytul: string): Promise<boolean>;
  /** Treść fragmentu wraz z jego postacią i blokadami, które go obejmują. */
  odczytajTekst(): Promise<string | null>;
  /** Nowe brzmienie zaznaczonego fragmentu bez przepisywania całości. */
  zmienTekst(brzmienie: string, zachowajPostac: boolean): Promise<boolean>;
  /** Pochodzenie fragmentów dokumentu — skąd każdy wniesiony fragment jest. */
  odczytajPochodzenia(): Promise<readonly StudioProvenance[]>;

  /* ── Nastawy widoku prowadzone przez rdzeń ───────────────────────────────── */

  /** Magazyn nastaw widoku oparty na `studio.view.get` i `.set`. */
  magazynWidoku: MagazynWidokuRdzenia;

  /** Null znaczy nieodczytaną jeszcze postać — zapis ma wtedy odczytać ją sam, nie wysłać postać pustą. */
  postacZapamietana(): StudioDocumentForm | null;
}

export function utworzPostacDokumentu(zaplecze: ZapleczePostaci): PostacDokumentu {
  const odpowiedz = utworzWierszOdpowiedzi();
  /** Uchwyt postaci pobranej malarzem formatów; pusty znaczy „nic nie pobrano". */
  let uchwytMalarza = '';

  const magazynWidoku = utworzMagazynWidokuRdzenia({
    zrodlo: zaplecze.zrodlo,
    idDokumentu: () => zaplecze.idDokumentu(),
    idOkna: () => zaplecze.idOkna(),
    naZdanie: (tresc, powodzenie) => zaplecze.naZdanie(tresc, powodzenie),
  });

  /* ── Wspólne dokładki żądań ──────────────────────────────────────────────── */

  /** Dokument czynny albo `null` wraz ze zdaniem o tym, czego brakuje. */
  function dokument(gdzie: (tresc: string, powodzenie: boolean) => void): string | null {
    const identyfikator = zaplecze.idDokumentu();
    if (identyfikator !== '') return identyfikator;
    gdzie(
      'Nie ma wczytanego dokumentu, a każda z tych komend przyjmuje jego identyfikator jako pole ' +
        'obowiązkowe. Wczytaj dokument albo załóż go z szablonu — brak jest po stronie okna, ' +
        'nie rdzenia.',
      false,
    );
    return null;
  }

  /** Zakres z zaznaczenia; brak zaznaczenia znaczy cały dokument, czyli brak pól. */
  function zakres(): { rangeStart?: number; rangeEnd?: number } {
    const wybrany = zaplecze.zaznaczenie();
    return wybrany === null
      ? {}
      : { rangeStart: wybrany.poczatek, rangeEnd: wybrany.koniec };
  }

  /** Autor czynności — panel należy do Operatora i nie udaje modelu. */
  const autorOperatora = { author: StudioAuthor.Uzytkownik } as const;

  /** Bilans idzie do Operatora zawsze, także gdy applied jest zerem — to informacja, nie cisza. */
  function ogloszSkutek(
    czynnosc: string,
    bilans: StudioActionBalance,
    postac: StudioDocumentForm,
    gdzie: (tresc: string, powodzenie: boolean) => void,
  ): void {
    zaplecze.naPostac(postac);
    const powodzenie = bilans.applied > 0 || bilans.skippedCount === 0;
    gdzie(`${czynnosc}: ${opiszBilans(bilans)}`, powodzenie);
  }

  /** Zdanie o odmowie rdzenia. */
  function ogloszOdmowe(
    czynnosc: string,
    blad: { code?: string; message?: string } | undefined,
    gdzie: (tresc: string, powodzenie: boolean) => void,
  ): void {
    gdzie(opisOdmowy(czynnosc, blad?.code, blad?.message), false);
  }

  /* ── Panel nastaw strony ─────────────────────────────────────────────────── */

  const panelStrony: StronaPanelNastaw = utworzStronePanelNastaw(
    {
      naNastawyStrony: (zadanie) => {
        const identyfikator = dokument(panelStrony.pokazOdpowiedz);
        if (identyfikator === null) return;
        void (async () => {
          const wynik = await zaplecze.zrodlo.ustawStrone({
            documentId: identyfikator,
            ...zadanie,
            ...autorOperatora,
          });
          if (!wynik.udany || wynik.wynik === undefined) {
            ogloszOdmowe('Nastawy strony', wynik.blad, panelStrony.pokazOdpowiedz);
            return;
          }
          panelStrony.pokazNastawy(wynik.wynik.pageSetup);
          ogloszSkutek(
            'Nastawy strony zapisane w rdzeniu',
            wynik.wynik.balance,
            wynik.wynik.form,
            panelStrony.pokazOdpowiedz,
          );
        })();
      },

      naNumeracje: (zadanie) => {
        const identyfikator = dokument(panelStrony.pokazOdpowiedz);
        if (identyfikator === null) return;
        void (async () => {
          const wynik = await zaplecze.zrodlo.ustawNumeracje({
            documentId: identyfikator,
            ...zadanie,
            ...autorOperatora,
          });
          if (!wynik.udany || wynik.wynik === undefined) {
            ogloszOdmowe('Numeracja stron', wynik.blad, panelStrony.pokazOdpowiedz);
            return;
          }
          const numeracja = wynik.wynik.numbering;
          ogloszSkutek(
            `Numeracja stron ${numeracja.enabled ? 'włączona' : 'wyłączona'}` +
              `${numeracja.format === undefined ? '' : `, styl ${numeracja.format}`}` +
              `${numeracja.startAt === undefined ? '' : `, od numeru ${numeracja.startAt}`}` +
              `${numeracja.restartInSection === true ? ', wznowiona w tej sekcji' : ''}` +
              `${numeracja.showTotal === true ? ', z liczbą stron' : ''}` +
              `${numeracja.position === undefined ? '' : `, umiejscowienie ${numeracja.position}`}`,
            wynik.wynik.balance,
            wynik.wynik.form,
            panelStrony.pokazOdpowiedz,
          );
        })();
      },

      naNaglowek: (zadanie) => {
        const identyfikator = dokument(panelStrony.pokazOdpowiedz);
        if (identyfikator === null) return;
        void (async () => {
          const wynik = await zaplecze.zrodlo.ustawNaglowek({
            documentId: identyfikator,
            ...zadanie,
            ...autorOperatora,
          });
          if (!wynik.udany || wynik.wynik === undefined) {
            ogloszOdmowe('Nagłówek i stopka', wynik.blad, panelStrony.pokazOdpowiedz);
            return;
          }
          panelStrony.pokazNaglowki(wynik.wynik.headersFooters);
          ogloszSkutek(
            `Nagłówek i stopka zapisane dla zasięgu ${zadanie.scope}`,
            wynik.wynik.balance,
            wynik.wynik.form,
            panelStrony.pokazOdpowiedz,
          );
        })();
      },

      naZnakWodny: (zadanie) => {
        const identyfikator = dokument(panelStrony.pokazOdpowiedz);
        if (identyfikator === null) return;
        void (async () => {
          const wynik = await zaplecze.zrodlo.ustawZnakWodny({
            documentId: identyfikator,
            ...zadanie,
            ...autorOperatora,
          });
          if (!wynik.udany || wynik.wynik === undefined) {
            ogloszOdmowe('Znak wodny', wynik.blad, panelStrony.pokazOdpowiedz);
            return;
          }
          ogloszSkutek(
            `Znak wodny zapisany — rodzaj ${wynik.wynik.watermark.kind}`,
            wynik.wynik.balance,
            wynik.wynik.form,
            panelStrony.pokazOdpowiedz,
          );
        })();
      },

      naKoperte: (zadanie) => {
        const identyfikator = dokument(panelStrony.pokazOdpowiedz);
        if (identyfikator === null) return;
        void (async () => {
          const wynik = await zaplecze.zrodlo.ustawKoperte({
            documentId: identyfikator,
            ...zadanie,
            ...autorOperatora,
          });
          if (!wynik.udany || wynik.wynik === undefined) {
            ogloszOdmowe('Nadruk koperty', wynik.blad, panelStrony.pokazOdpowiedz);
            return;
          }
          ogloszSkutek(
            'Nadruk koperty zapisany',
            wynik.wynik.balance,
            wynik.wynik.form,
            panelStrony.pokazOdpowiedz,
          );
        })();
      },

      naPodzial: (zadanie) => {
        const identyfikator = dokument(panelStrony.pokazOdpowiedz);
        if (identyfikator === null) return;
        void (async () => {
          const wynik = await zaplecze.zrodlo.wstawPodzial({
            documentId: identyfikator,
            ...zadanie,
            ...autorOperatora,
          });
          if (!wynik.udany || wynik.wynik === undefined) {
            ogloszOdmowe('Wstawienie podziału', wynik.blad, panelStrony.pokazOdpowiedz);
            return;
          }
          ogloszSkutek(
            `Podział „${zadanie.kind}" wstawiony na znaku ${zadanie.offset}`,
            wynik.wynik.balance,
            wynik.wynik.form,
            panelStrony.pokazOdpowiedz,
          );
        })();
      },

      naZapisSekcji: (zadanie) => {
        const identyfikator = dokument(panelStrony.pokazOdpowiedz);
        if (identyfikator === null) return;
        void (async () => {
          const wynik = await zaplecze.zrodlo.zapiszSekcje({
            documentId: identyfikator,
            ...zadanie,
            ...autorOperatora,
          });
          if (!wynik.udany || wynik.wynik === undefined) {
            ogloszOdmowe('Zapis sekcji', wynik.blad, panelStrony.pokazOdpowiedz);
            return;
          }
          ogloszSkutek(
            `Sekcja ${wynik.wynik.section.id} zapisana`,
            wynik.wynik.balance,
            wynik.wynik.form,
            panelStrony.pokazOdpowiedz,
          );
          void odczytajSekcje(identyfikator);
        })();
      },

      naUsuniecieSekcji: (idSekcji) => {
        const identyfikator = dokument(panelStrony.pokazOdpowiedz);
        if (identyfikator === null) return;
        void (async () => {
          const wynik = await zaplecze.zrodlo.usunSekcje({
            documentId: identyfikator,
            sectionId: idSekcji,
            ...autorOperatora,
          });
          if (!wynik.udany || wynik.wynik === undefined) {
            ogloszOdmowe('Usunięcie sekcji', wynik.blad, panelStrony.pokazOdpowiedz);
            return;
          }
          if (!wynik.wynik.deleted) {
            // Deleted: false z bilansem jest odpowiedzią poprawną — rdzeń odmówił i powiedział, dlaczego.
            panelStrony.pokazOdpowiedz(
              `Rdzeń NIE usunął sekcji ${idSekcji}. ${opiszBilans(wynik.wynik.balance)}`,
              false,
            );
            return;
          }
          ogloszSkutek(
            `Sekcja ${idSekcji} usunięta — jej treść przeszła do sekcji poprzedniej`,
            wynik.wynik.balance,
            wynik.wynik.form,
            panelStrony.pokazOdpowiedz,
          );
          void odczytajSekcje(identyfikator);
        })();
      },

      naOdczyt: () => void odczytajStrone(),
    },
    () => zaplecze.miejsceKursora(),
  );

  /* ── Panel arkusza stylów ────────────────────────────────────────────────── */

  const panelStylu: StylPanelArkusza = utworzStylPanelArkusza({
    naZapisStylu: (zadanie) => {
      const identyfikator = dokument(panelStylu.pokazOdpowiedz);
      if (identyfikator === null) return;
      void (async () => {
        const wynik = await zaplecze.zrodlo.zapiszStyl({
          documentId: identyfikator,
          ...zadanie,
          ...autorOperatora,
        });
        if (!wynik.udany || wynik.wynik === undefined) {
          ogloszOdmowe('Zapis stylu nazwanego', wynik.blad, panelStylu.pokazOdpowiedz);
          return;
        }
        ogloszSkutek(
          `Styl „${wynik.wynik.style.name}" zapisany; miejsc użycia ` +
            `${wynik.wynik.style.usageCount ?? 0} — zmiana stylu przestawia je wszystkie`,
          wynik.wynik.balance,
          wynik.wynik.form,
          panelStylu.pokazOdpowiedz,
        );
        void odczytajStyle(identyfikator);
      })();
    },

    naStosowanieStylu: (zadanie) => {
      const identyfikator = dokument(panelStylu.pokazOdpowiedz);
      if (identyfikator === null) return;
      void (async () => {
        const wynik = await zaplecze.zrodlo.zastosujStyl({
          documentId: identyfikator,
          ...zadanie,
          ...zakres(),
          ...autorOperatora,
        });
        if (!wynik.udany || wynik.wynik === undefined) {
          ogloszOdmowe('Zastosowanie stylu', wynik.blad, panelStylu.pokazOdpowiedz);
          return;
        }
        ogloszSkutek(
          `Styl „${zadanie.name}" zastosowany`,
          wynik.wynik.balance,
          wynik.wynik.form,
          panelStylu.pokazOdpowiedz,
        );
      })();
    },

    naUsuniecieStylu: (zadanie) => {
      const identyfikator = dokument(panelStylu.pokazOdpowiedz);
      if (identyfikator === null) return;
      void (async () => {
        const wynik = await zaplecze.zrodlo.usunStyl({
          documentId: identyfikator,
          ...zadanie,
          ...autorOperatora,
        });
        if (!wynik.udany || wynik.wynik === undefined) {
          ogloszOdmowe('Usunięcie stylu', wynik.blad, panelStylu.pokazOdpowiedz);
          return;
        }
        if (!wynik.wynik.deleted) {
          panelStylu.pokazOdpowiedz(
            `Rdzeń NIE usunął stylu „${zadanie.name}" — tak odpowiada dla stylu fabrycznego. ` +
              opiszBilans(wynik.wynik.balance),
            false,
          );
          return;
        }
        ogloszSkutek(
          `Styl „${zadanie.name}" usunięty`,
          wynik.wynik.balance,
          wynik.wynik.form,
          panelStylu.pokazOdpowiedz,
        );
        void odczytajStyle(identyfikator);
      })();
    },

    naStylZnaku: (zadanie) => {
      const identyfikator = dokument(panelStylu.pokazOdpowiedz);
      if (identyfikator === null) return;
      void (async () => {
        const wynik = await zaplecze.zrodlo.ustawZnak({
          documentId: identyfikator,
          ...zadanie,
          ...zakres(),
          ...autorOperatora,
        });
        if (!wynik.udany || wynik.wynik === undefined) {
          ogloszOdmowe('Styl znaku', wynik.blad, panelStylu.pokazOdpowiedz);
          return;
        }
        ogloszSkutek(
          'Styl znaku naniesiony',
          wynik.wynik.balance,
          wynik.wynik.form,
          panelStylu.pokazOdpowiedz,
        );
      })();
    },

    naStylAkapitu: (zadanie) => void nanieStylAkapitu(zadanie, panelStylu.pokazOdpowiedz),

    naCzyszczenie: (zadanie) => {
      const identyfikator = dokument(panelStylu.pokazOdpowiedz);
      if (identyfikator === null) return;
      void (async () => {
        const wynik = await zaplecze.zrodlo.wyczyscFormat({
          documentId: identyfikator,
          ...zadanie,
          ...zakres(),
          ...autorOperatora,
        });
        if (!wynik.udany || wynik.wynik === undefined) {
          ogloszOdmowe('Czyszczenie formatowania', wynik.blad, panelStylu.pokazOdpowiedz);
          return;
        }
        ogloszSkutek(
          'Formatowanie wyczyszczone — litery zostały',
          wynik.wynik.balance,
          wynik.wynik.form,
          panelStylu.pokazOdpowiedz,
        );
      })();
    },

    naWielkoscLiter: (zadanie) => {
      const identyfikator = dokument(panelStylu.pokazOdpowiedz);
      if (identyfikator === null) return;
      void (async () => {
        const wynik = await zaplecze.zrodlo.ustawWielkoscLiter({
          documentId: identyfikator,
          ...zadanie,
          ...zakres(),
          ...autorOperatora,
        });
        if (!wynik.udany || wynik.wynik === undefined) {
          ogloszOdmowe('Wielkość liter', wynik.blad, panelStylu.pokazOdpowiedz);
          return;
        }
        ogloszSkutek(
          `Wielkość liter przestawiona (${zadanie.transform})`,
          wynik.wynik.balance,
          wynik.wynik.form,
          panelStylu.pokazOdpowiedz,
        );
      })();
    },

    naPobraniePostaci: (zadanie) => {
      const identyfikator = dokument(panelStylu.pokazOdpowiedz);
      if (identyfikator === null) return;
      void (async () => {
        const wynik = await zaplecze.zrodlo.pobierzPostac({
          documentId: identyfikator,
          ...zadanie,
          ...zakres(),
        });
        if (!wynik.udany || wynik.wynik === undefined) {
          ogloszOdmowe('Malarz formatów — pobranie', wynik.blad, panelStylu.pokazOdpowiedz);
          return;
        }
        uchwytMalarza = wynik.wynik.clipId;
        const znak = wynik.wynik.character;
        const akapit = wynik.wynik.paragraph;
        panelStylu.pokazMalarza(
          true,
          `uchwyt ${uchwytMalarza}` +
            `${znak?.fontFamily === undefined ? '' : `, krój ${znak.fontFamily}`}` +
            `${znak?.fontSizePt === undefined ? '' : `, stopień ${znak.fontSizePt} pt`}` +
            `${akapit?.align === undefined ? '' : `, wyrównanie ${akapit.align}`}`,
        );
        panelStylu.pokazOdpowiedz(
          'Postać pobrana. Zaznacz miejsce docelowe i naciśnij „Malarz: nanieś na zaznaczenie" — ' +
            'malarz kopiuje POSTAĆ, nie treść.',
          true,
        );
      })();
    },

    naNalozeniePostaci: (zadanie) => {
      const identyfikator = dokument(panelStylu.pokazOdpowiedz);
      if (identyfikator === null) return;
      if (uchwytMalarza === '') {
        panelStylu.pokazOdpowiedz(
          'Malarz formatów nie ma nic pobranego — najpierw pobierz postać ze wzorcowego fragmentu. ' +
            'Bez uchwytu rdzeń nie ma czego nanieść.',
          false,
        );
        return;
      }
      void (async () => {
        const wynik = await zaplecze.zrodlo.nalozPostac({
          documentId: identyfikator,
          clipId: uchwytMalarza,
          ...zadanie,
          ...zakres(),
          ...autorOperatora,
        });
        if (!wynik.udany || wynik.wynik === undefined) {
          ogloszOdmowe('Malarz formatów — naniesienie', wynik.blad, panelStylu.pokazOdpowiedz);
          return;
        }
        ogloszSkutek(
          'Pobrana postać naniesiona',
          wynik.wynik.balance,
          wynik.wynik.form,
          panelStylu.pokazOdpowiedz,
        );
      })();
    },

    naPodobne: (zadanie) => {
      const identyfikator = dokument(panelStylu.pokazOdpowiedz);
      if (identyfikator === null) return;
      void (async () => {
        const wynik = await zaplecze.zrodlo.podobnePostacia({
          documentId: identyfikator,
          ...zadanie,
          ...zakres(),
        });
        if (!wynik.udany || wynik.wynik === undefined) {
          ogloszOdmowe('Zaznaczenie wedle postaci', wynik.blad, panelStylu.pokazOdpowiedz);
          return;
        }
        const trafienia = wynik.wynik.matches;
        panelStylu.pokazOdpowiedz(
          wynik.wynik.count === 0
            ? 'Rdzeń nie znalazł ani jednego fragmentu o postaci zgodnej ze wzorem.'
            : `Fragmentów o zgodnej postaci: ${wynik.wynik.count}. Pierwszy stoi na znakach ` +
                `${trafienia[0]?.rangeStart ?? 0}–${trafienia[0]?.rangeEnd ?? 0}. Wykaz oddaje rdzeń ` +
                'zakresami — okno przenosi na nie zaznaczenie po jednym, bo zaznaczenie w tym ' +
                'edytorze jest jedno, a nie wielokrotne.',
          wynik.wynik.count > 0,
        );
      })();
    },

    naZamiane: (zadanie) => {
      const identyfikator = dokument(panelStylu.pokazOdpowiedz);
      if (identyfikator === null) return;
      void (async () => {
        const wynik = await zaplecze.zrodlo.zamienZPostacia({
          documentId: identyfikator,
          ...zadanie,
          ...(panelStylu.zamianaTylkoWZaznaczeniu() ? zakres() : {}),
          ...autorOperatora,
        });
        if (!wynik.udany || wynik.wynik === undefined) {
          ogloszOdmowe('Zamiana wraz z postacią', wynik.blad, panelStylu.pokazOdpowiedz);
          return;
        }
        ogloszSkutek(
          `Trafień ${wynik.wynik.matches}, zamienionych ${wynik.wynik.replaced}`,
          wynik.wynik.balance,
          wynik.wynik.form,
          panelStylu.pokazOdpowiedz,
        );
      })();
    },

    naTabulator: (zadanie) => {
      const identyfikator = dokument(panelStylu.pokazOdpowiedz);
      if (identyfikator === null) return;
      void (async () => {
        const wynik = await zaplecze.zrodlo.ustawTabulator({
          documentId: identyfikator,
          ...zadanie,
          ...zakres(),
          ...autorOperatora,
        });
        if (!wynik.udany || wynik.wynik === undefined) {
          ogloszOdmowe('Tabulator', wynik.blad, panelStylu.pokazOdpowiedz);
          return;
        }
        ogloszSkutek(
          `Tabulatory akapitu po zmianie: ${wynik.wynik.tabStops.length}`,
          wynik.wynik.balance,
          wynik.wynik.form,
          panelStylu.pokazOdpowiedz,
        );
      })();
    },

    naOdczyt: () => void odczytajStyleIPostac(),
  });

  /* ── Panel list i znaków ─────────────────────────────────────────────────── */

  const panelList: StylPanelList = utworzStylPanelList({
    naListe: (zadanie) => {
      const identyfikator = dokument(panelList.pokazOdpowiedz);
      if (identyfikator === null) return;
      void (async () => {
        const wynik = await zaplecze.zrodlo.zastosujListe({
          documentId: identyfikator,
          ...zadanie,
          ...zakres(),
          ...autorOperatora,
        });
        if (!wynik.udany || wynik.wynik === undefined) {
          ogloszOdmowe('Lista', wynik.blad, panelList.pokazOdpowiedz);
          return;
        }
        ogloszSkutek(
          `Lista ${wynik.wynik.list.id} (${wynik.wynik.list.kind})`,
          wynik.wynik.balance,
          wynik.wynik.form,
          panelList.pokazOdpowiedz,
        );
        panelList.pokazListy(wynik.wynik.form.lists ?? []);
      })();
    },

    naPunktator: (zadanie) => {
      const identyfikator = dokument(panelList.pokazOdpowiedz);
      if (identyfikator === null) return;
      void (async () => {
        const wynik = await zaplecze.zrodlo.ustawPunktator({
          documentId: identyfikator,
          ...zadanie,
          ...autorOperatora,
        });
        if (!wynik.udany || wynik.wynik === undefined) {
          ogloszOdmowe('Znak wypunktowania', wynik.blad, panelList.pokazOdpowiedz);
          return;
        }
        ogloszSkutek(
          `Znak wypunktowania poziomu ${zadanie.level} zapisany`,
          wynik.wynik.balance,
          wynik.wynik.form,
          panelList.pokazOdpowiedz,
        );
      })();
    },

    naNumeracjeListy: (zadanie) => {
      const identyfikator = dokument(panelList.pokazOdpowiedz);
      if (identyfikator === null) return;
      void (async () => {
        const wynik = await zaplecze.zrodlo.ustawNumeracjeListy({
          documentId: identyfikator,
          ...zadanie,
          ...autorOperatora,
        });
        if (!wynik.udany || wynik.wynik === undefined) {
          ogloszOdmowe('Numeracja poziomu listy', wynik.blad, panelList.pokazOdpowiedz);
          return;
        }
        ogloszSkutek(
          `Format numeracji poziomu ${zadanie.level} zapisany`,
          wynik.wynik.balance,
          wynik.wynik.form,
          panelList.pokazOdpowiedz,
        );
      })();
    },

    naWznowienie: (zadanie) => {
      const identyfikator = dokument(panelList.pokazOdpowiedz);
      if (identyfikator === null) return;
      void (async () => {
        const wynik = await zaplecze.zrodlo.wznowNumeracje({
          documentId: identyfikator,
          ...zadanie,
          // Miejsce wznowienia bierze się z kursora, nie z pola — wznawia się tam, gdzie stoi Operator.
          offset: zaplecze.miejsceKursora(),
          ...autorOperatora,
        });
        if (!wynik.udany || wynik.wynik === undefined) {
          ogloszOdmowe('Wznowienie numeracji', wynik.blad, panelList.pokazOdpowiedz);
          return;
        }
        ogloszSkutek(
          `Numeracja wznowiona na znaku ${zaplecze.miejsceKursora()}`,
          wynik.wynik.balance,
          wynik.wynik.form,
          panelList.pokazOdpowiedz,
        );
      })();
    },

    naPoziom: (zadanie) => {
      const identyfikator = dokument(panelList.pokazOdpowiedz);
      if (identyfikator === null) return;
      void (async () => {
        const wynik = await zaplecze.zrodlo.przestawPoziom({
          documentId: identyfikator,
          ...zadanie,
          ...zakres(),
          ...autorOperatora,
        });
        if (!wynik.udany || wynik.wynik === undefined) {
          ogloszOdmowe('Poziom listy', wynik.blad, panelList.pokazOdpowiedz);
          return;
        }
        ogloszSkutek(
          `Poziom listy przestawiony o ${zadanie.step}`,
          wynik.wynik.balance,
          wynik.wynik.form,
          panelList.pokazOdpowiedz,
        );
      })();
    },

    naZnak: (zadanie) => {
      const identyfikator = dokument(panelList.pokazOdpowiedz);
      if (identyfikator === null) return;
      void (async () => {
        const wynik = await zaplecze.zrodlo.wstawZnak({
          documentId: identyfikator,
          offset: zaplecze.miejsceKursora(),
          ...zadanie,
          ...autorOperatora,
        });
        if (!wynik.udany || wynik.wynik === undefined) {
          ogloszOdmowe('Wstawienie znaku', wynik.blad, panelList.pokazOdpowiedz);
          return;
        }
        ogloszSkutek(
          `Znak „${wynik.wynik.symbol.character}" (${wynik.wynik.symbol.name}) wstawiony`,
          wynik.wynik.balance,
          wynik.wynik.form,
          panelList.pokazOdpowiedz,
        );
        // Wykaz ostatnio użytych znaków odczytuje się po wstawieniu, żeby znak trafił pod rękę od razu.
        void odczytajZnaki('', false);
      })();
    },

    naSzukanieZnaku: (fraza, tylkoOstatnie) => void odczytajZnaki(fraza, tylkoOstatnie),

    naAutozamiane: (zadanie) => {
      void (async () => {
        const wynik = await zaplecze.zrodlo.ustawAutozamiane(zadanie);
        if (!wynik.udany || wynik.wynik === undefined) {
          ogloszOdmowe('Zasada autozamiany', wynik.blad, panelList.pokazOdpowiedz);
          return;
        }
        panelList.pokazOdpowiedz(
          wynik.wynik.removed === true
            ? `Zasada autozamiany „${zadanie.shortcut}" usunięta.`
            : `Zasada autozamiany „${zadanie.shortcut}" zapisana.`,
          true,
        );
        void odczytajAutozamiany();
      })();
    },

    naOdczyt: () => {
      void odczytajZnaki('', false);
      void odczytajAutozamiany();
      const identyfikator = zaplecze.idDokumentu();
      if (identyfikator !== '') void odczytajListy(identyfikator);
    },
  });

  /* ── Panel treści, postaci i pochodzenia ─────────────────────────────────── */

  const panelTresci: StronaPanelTresci = utworzStronePanelTresci({
    naOdczytFragmentu: () => {
      const identyfikator = dokument(panelTresci.pokazOdpowiedz);
      if (identyfikator === null) return;
      void (async () => {
        const wynik = await zaplecze.zrodlo.tekst({ documentId: identyfikator, ...zakres() });
        if (!wynik.udany || wynik.wynik === undefined) {
          ogloszOdmowe('Odczyt fragmentu', wynik.blad, panelTresci.pokazOdpowiedz);
          return;
        }
        const blokady = wynik.wynik.locks ?? [];
        panelTresci.pokazFragment(
          wynik.wynik.text,
          blokady.length === 0
            ? 'Fragmentu nie obejmuje żadna blokada.'
            : `Fragment obejmują blokady: ${blokady
                .map(
                  (blokada) =>
                    `„${blokada.name}" (znaki ${blokada.rangeStart}–${blokada.rangeEnd}, zasięg ` +
                    `${blokada.scope})`,
                )
                .join(', ')}. Sprawdzenie stoi w RDZENIU, przed dotknięciem treści.`,
        );
        panelTresci.pokazOdpowiedz(
          `Odczytano ${wynik.wynik.text.length} znaków; fragmentów o jednolitej postaci ` +
            `${wynik.wynik.runs?.length ?? 0}.`,
          true,
        );
      })();
    },

    naZmianeBrzmienia: (nowe, zachowaj) => {
      const identyfikator = dokument(panelTresci.pokazOdpowiedz);
      if (identyfikator === null) return;
      const wybrany = zaplecze.zaznaczenie();
      if (wybrany === null) {
        panelTresci.pokazOdpowiedz(
          'Zmiana brzmienia dotyczy FRAGMENTU — zaznacz go. Komenda studio.text.edit przyjmuje ' +
            'początek i koniec jako pola obowiązkowe, bo jej sensem jest poprawa bez przepisywania ' +
            'całości.',
          false,
        );
        return;
      }
      void (async () => {
        const wynik = await zaplecze.zrodlo.zmienTekst({
          documentId: identyfikator,
          rangeStart: wybrany.poczatek,
          rangeEnd: wybrany.koniec,
          text: nowe,
          keepFormat: zachowaj,
          ...autorOperatora,
        });
        if (!wynik.udany || wynik.wynik === undefined) {
          ogloszOdmowe('Zmiana brzmienia fragmentu', wynik.blad, panelTresci.pokazOdpowiedz);
          return;
        }
        ogloszSkutek(
          `Brzmienie fragmentu (znaki ${wybrany.poczatek}–${wybrany.koniec}) zmienione`,
          wynik.wynik.balance,
          wynik.wynik.form,
          panelTresci.pokazOdpowiedz,
        );
      })();
    },

    naOdczytPostaci: () => {
      const identyfikator = dokument(panelTresci.pokazOdpowiedz);
      if (identyfikator === null) return;
      void (async () => {
        const wynik = await zaplecze.zrodlo.postac({ documentId: identyfikator });
        if (!wynik.udany || wynik.wynik === undefined) {
          ogloszOdmowe('Odczyt postaci dokumentu', wynik.blad, panelTresci.pokazOdpowiedz);
          return;
        }
        postacOstatnia = wynik.wynik.form;
        zaplecze.naPostac(wynik.wynik.form);
        panelTresci.pokazPostac(opiszPostac(wynik.wynik.form));
        panelTresci.pokazOdpowiedz('Postać dokumentu odczytana z rdzenia.', true);
      })();
    },

    naZapisPostaci: () => {
      const identyfikator = dokument(panelTresci.pokazOdpowiedz);
      if (identyfikator === null) return;
      void (async () => {
        // Zapisuje się postać odczytaną teraz, nie zapamiętaną wcześniej — inne okno mogło ją zmienić.
        const odczyt = await zaplecze.zrodlo.postac({ documentId: identyfikator });
        if (!odczyt.udany || odczyt.wynik === undefined) {
          ogloszOdmowe('Odczyt postaci przed zapisem', odczyt.blad, panelTresci.pokazOdpowiedz);
          return;
        }
        const wynik = await zaplecze.zrodlo.zapiszPostac({
          documentId: identyfikator,
          form: odczyt.wynik.form,
          // Treści nie ma w żądaniu: brak pola znaczy bez zmiany treści, treść pusta skasowałaby dokument.
          createVersion: true,
          ...autorOperatora,
        });
        if (!wynik.udany || wynik.wynik === undefined) {
          ogloszOdmowe('Utrwalenie postaci', wynik.blad, panelTresci.pokazOdpowiedz);
          return;
        }
        postacOstatnia = wynik.wynik.form;
        ogloszSkutek(
          'Postać utrwalona bez dotykania treści' +
            `${wynik.wynik.version === undefined ? '' : `; wersja ${wynik.wynik.version.id}`}`,
          wynik.wynik.balance,
          wynik.wynik.form,
          panelTresci.pokazOdpowiedz,
        );
      })();
    },

    naPochodzenia: () => {
      const identyfikator = dokument(panelTresci.pokazOdpowiedz);
      if (identyfikator === null) return;
      void (async () => {
        const wynik = await zaplecze.zrodlo.pochodzenia({ documentId: identyfikator });
        if (!wynik.udany || wynik.wynik === undefined) {
          ogloszOdmowe('Pochodzenie fragmentów', wynik.blad, panelTresci.pokazOdpowiedz);
          return;
        }
        panelTresci.pokazPochodzenia(wynik.wynik.entries);
        panelTresci.pokazOdpowiedz(
          `Zapisów pochodzenia: ${wynik.wynik.entries.length}.`,
          true,
        );
      })();
    },
  });

  /** Postać odczytana ostatnio — czyta ją okno przez `postacZapamietana`. */
  let postacOstatnia: StudioDocumentForm | null = null;

  /* ── Odczyty ─────────────────────────────────────────────────────────────── */

  async function odczytajStrone(): Promise<void> {
    const identyfikator = dokument(panelStrony.pokazOdpowiedz);
    void odczytajNosniki();
    if (identyfikator === null) return;
    const sekcja = panelStrony.sekcjaWybrana();
    const nastawy = await zaplecze.zrodlo.nastawyStrony({
      documentId: identyfikator,
      ...(sekcja === undefined ? {} : { sectionId: sekcja }),
    });
    if (!nastawy.udany || nastawy.wynik === undefined) {
      ogloszOdmowe('Odczyt nastaw strony', nastawy.blad, panelStrony.pokazOdpowiedz);
      return;
    }
    panelStrony.pokazNastawy(nastawy.wynik.pageSetup);
    if (nastawy.wynik.sections !== undefined) panelStrony.pokazSekcje(nastawy.wynik.sections);
    else await odczytajSekcje(identyfikator);
    const naglowki = await zaplecze.zrodlo.naglowki({
      documentId: identyfikator,
      ...(sekcja === undefined ? {} : { sectionId: sekcja }),
    });
    if (naglowki.udany && naglowki.wynik !== undefined) {
      panelStrony.pokazNaglowki(naglowki.wynik.headersFooters);
    }
    panelStrony.pokazOdpowiedz(
      'Pola pokazują nastawy odczytane z rdzenia — to, co naprawdę stoi w dokumencie.',
      true,
    );
  }

  async function odczytajNosniki(): Promise<void> {
    const wynik = await zaplecze.zrodlo.nosniki({});
    if (!wynik.udany || wynik.wynik === undefined) return;
    panelStrony.pokazNosniki(wynik.wynik.papers);
    // Wykaz rdzenia przestawia też wymiary kartki, nie tylko nazwy nośników do wyboru.
    wchlonNosnikiKontraktu(wynik.wynik.papers);
  }

  async function odczytajSekcje(identyfikator: string): Promise<void> {
    const wynik = await zaplecze.zrodlo.sekcje({ documentId: identyfikator });
    if (!wynik.udany || wynik.wynik === undefined) return;
    panelStrony.pokazSekcje(wynik.wynik.sections);
  }

  async function odczytajStyle(identyfikator: string): Promise<void> {
    const wynik = await zaplecze.zrodlo.style({ documentId: identyfikator });
    if (!wynik.udany || wynik.wynik === undefined) {
      ogloszOdmowe('Odczyt arkusza stylów', wynik.blad, panelStylu.pokazOdpowiedz);
      return;
    }
    panelStylu.pokazStyle(wynik.wynik.styles);
  }

  async function odczytajStyleIPostac(): Promise<void> {
    const identyfikator = dokument(panelStylu.pokazOdpowiedz);
    if (identyfikator === null) return;
    await odczytajStyle(identyfikator);
    const znak = await zaplecze.zrodlo.postacZnaku({ documentId: identyfikator, ...zakres() });
    if (znak.udany && znak.wynik !== undefined) {
      panelStylu.pokazPostacZnaku(
        znak.wynik.character as StudioCharacterFormat,
        znak.wynik.mixedFields ?? [],
      );
    }
    const akapit = await zaplecze.zrodlo.postacAkapitu({ documentId: identyfikator, ...zakres() });
    if (akapit.udany && akapit.wynik !== undefined) {
      panelStylu.pokazPostacAkapitu(
        akapit.wynik.paragraph as StudioParagraphFormat,
        akapit.wynik.mixedFields ?? [],
      );
    }
  }

  async function odczytajListy(identyfikator: string): Promise<void> {
    const wynik = await zaplecze.zrodlo.postac({
      documentId: identyfikator,
      // Bloki treści nie są tu potrzebne — wykaz list ma osobne pole, bloki byłyby setkami pozycji.
      includeBlocks: false,
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      ogloszOdmowe('Odczyt definicji list', wynik.blad, panelList.pokazOdpowiedz);
      return;
    }
    panelList.pokazListy(wynik.wynik.form.lists ?? []);
  }

  async function odczytajZnaki(fraza: string, tylkoOstatnie: boolean): Promise<void> {
    const wynik = await zaplecze.zrodlo.znaki({
      ...(fraza === '' ? {} : { query: fraza }),
      ...(tylkoOstatnie ? { recentOnly: true } : {}),
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      ogloszOdmowe('Tablica znaków', wynik.blad, panelList.pokazOdpowiedz);
      return;
    }
    panelList.pokazZnaki(wynik.wynik.symbols, wynik.wynik.categories ?? []);
  }

  async function odczytajAutozamiany(): Promise<void> {
    const wynik = await zaplecze.zrodlo.autozamiany({});
    if (!wynik.udany || wynik.wynik === undefined) return;
    panelList.pokazAutozamiany(wynik.wynik.rules);
  }

  /* ── Styl akapitu wołany z dwóch stron ───────────────────────────────────── */

  /** Chwyt wcięcia na linijce i pole panelu ustawiają tę samą cechę akapitu tą samą komendą rdzenia. */
  async function nanieStylAkapitu(
    zadanie: Record<string, unknown>,
    gdzie: (tresc: string, powodzenie: boolean) => void,
  ): Promise<void> {
    const identyfikator = dokument(gdzie);
    if (identyfikator === null) return;
    const wynik = await zaplecze.zrodlo.ustawAkapit({
      documentId: identyfikator,
      ...zadanie,
      ...zakres(),
      ...autorOperatora,
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      ogloszOdmowe('Styl akapitu', wynik.blad, gdzie);
      return;
    }
    ogloszSkutek('Styl akapitu naniesiony', wynik.wynik.balance, wynik.wynik.form, gdzie);
  }

  /* ── Tabulatory: różnica wobec wykazu zapamiętanego ──────────────────────── */

  /** Linijka oddaje wykaz tabulatorów w całości, rdzeń przyjmuje po jednym — pamięć wykazu chroni resztę. */
  let tabulatoryAkapitu: TabulatorAkapitu[] = [];

  /** Rodzaj tabulatora linijki przełożony na kontrakt. */
  function rodzajKontraktu(rodzaj: RodzajTabulatora): StudioTabKind {
    if (rodzaj === 'prawy') return StudioTabKind.Right;
    if (rodzaj === 'srodkowy') return StudioTabKind.Center;
    if (rodzaj === 'dziesietny') return StudioTabKind.Decimal;
    return StudioTabKind.Left;
  }

  /** Znak wiodący linijki przełożony na kontrakt. */
  function znakKontraktu(znak: ZnakWiodacy): StudioTabLeader {
    if (znak === 'kropka') return StudioTabLeader.Dot;
    if (znak === 'kreska') return StudioTabLeader.Dash;
    if (znak === 'podkreslenie') return StudioTabLeader.Underline;
    return StudioTabLeader.None;
  }

  async function wyslijTabulator(
    identyfikator: string,
    tabulator: TabulatorAkapitu,
    zdejmij: boolean,
  ): Promise<void> {
    const wynik = await zaplecze.zrodlo.ustawTabulator({
      documentId: identyfikator,
      positionMm: tabulator.milimetry,
      kind: rodzajKontraktu(tabulator.rodzaj),
      leader: znakKontraktu(tabulator.znakWiodacy),
      remove: zdejmij,
      ...zakres(),
      ...autorOperatora,
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      ogloszOdmowe('Tabulator z linijki', wynik.blad, zaplecze.naZdanie);
      return;
    }
    ogloszSkutek(
      `Tabulator na ${tabulator.milimetry} mm ${zdejmij ? 'zdjęty' : 'zapisany'}; tabulatorów ` +
        `akapitu ${wynik.wynik.tabStops.length}`,
      wynik.wynik.balance,
      wynik.wynik.form,
      zaplecze.naZdanie,
    );
  }

  /* ── Pas otwarć ──────────────────────────────────────────────────────────── */

  const otworzStrone = przycisk('Nastawy strony i numeracja', 'dn-btn dn-btn--sm dn-btn--zarys');
  otworzStrone.dataset['czynnosc'] = 'otworz-nastawy-strony';
  otworzStrone.addEventListener('click', () => {
    panelStrony.przestawWidocznosc();
    if (panelStrony.widoczny()) void odczytajStrone();
  });

  const otworzStyl = przycisk('Arkusz stylów i formatowanie', 'dn-btn dn-btn--sm dn-btn--zarys');
  otworzStyl.dataset['czynnosc'] = 'otworz-arkusz-stylow';
  otworzStyl.addEventListener('click', () => {
    panelStylu.przestawWidocznosc();
    if (panelStylu.widoczny()) void odczytajStyleIPostac();
  });

  const otworzTresc = przycisk('Fragment, postać i pochodzenie', 'dn-btn dn-btn--sm dn-btn--zarys');
  otworzTresc.dataset['czynnosc'] = 'otworz-tresc-i-postac';
  otworzTresc.addEventListener('click', () => panelTresci.przestawWidocznosc());

  const otworzListy = przycisk('Listy i znaki specjalne', 'dn-btn dn-btn--sm dn-btn--zarys');
  otworzListy.dataset['czynnosc'] = 'otworz-listy';
  otworzListy.addEventListener('click', () => {
    panelList.przestawWidocznosc();
    if (!panelList.widoczny()) return;
    void odczytajZnaki('', false);
    void odczytajAutozamiany();
    const identyfikator = zaplecze.idDokumentu();
    if (identyfikator !== '') void odczytajListy(identyfikator);
  });

  const pas = document.createElement('div');
  pas.className = 'ms-postac__pas';
  pas.append(otworzStrone, otworzStyl, otworzListy, otworzTresc);

  const element = document.createElement('div');
  element.className = 'ms-postac__powloka';
  element.append(
    pas,
    panelStrony.element,
    panelStylu.element,
    panelList.element,
    panelTresci.element,
    odpowiedz.element,
  );

  panelStylu.pokazMalarza(false, '');

  return {
    element,

    async wczytaj() {
      await magazynWidoku.wczytaj();
      await odczytajNosniki();
      await odczytajZnaki('', false);
      await odczytajAutozamiany();
      const identyfikator = zaplecze.idDokumentu();
      if (identyfikator === '') return;
      await odczytajStrone();
      await odczytajStyle(identyfikator);
      await odczytajListy(identyfikator);
    },

    odswiez() {
      const identyfikator = zaplecze.idDokumentu();
      if (identyfikator === '') return;
      if (panelStrony.widoczny()) void odczytajStrone();
      if (panelStylu.widoczny()) void odczytajStyleIPostac();
      if (panelList.widoczny()) void odczytajListy(identyfikator);
    },

    zglosMargines(ktory, milimetry) {
      const identyfikator = dokument(zaplecze.naZdanie);
      if (identyfikator === null) return;
      const sekcja = panelStrony.sekcjaWybrana();
      const pole =
        ktory === 'gora'
          ? { marginTopMm: milimetry }
          : ktory === 'dol'
            ? { marginBottomMm: milimetry }
            : ktory === 'lewy'
              ? { marginLeftMm: milimetry }
              : { marginRightMm: milimetry };
      void (async () => {
        const wynik = await zaplecze.zrodlo.ustawStrone({
          documentId: identyfikator,
          ...(sekcja === undefined ? {} : { sectionId: sekcja }),
          ...pole,
          ...autorOperatora,
        });
        if (!wynik.udany || wynik.wynik === undefined) {
          ogloszOdmowe('Margines z linijki', wynik.blad, zaplecze.naZdanie);
          return;
        }
        panelStrony.pokazNastawy(wynik.wynik.pageSetup);
        ogloszSkutek(
          `Margines ${ktory} ustawiony na ${milimetry} mm i zapisany w rdzeniu`,
          wynik.wynik.balance,
          wynik.wynik.form,
          zaplecze.naZdanie,
        );
      })();
    },

    zglosStrone(strona) {
      const identyfikator = dokument(zaplecze.naZdanie);
      if (identyfikator === null) return;
      const sekcja = panelStrony.sekcjaWybrana();
      void (async () => {
        const wynik = await zaplecze.zrodlo.ustawStrone({
          documentId: identyfikator,
          ...(sekcja === undefined ? {} : { sectionId: sekcja }),
          paperName: strona.nosnik.oznaczenie,
          orientation: strona.orientacja,
          marginTopMm: strona.marginesGoraMm,
          marginBottomMm: strona.marginesDolMm,
          marginLeftMm: strona.marginesLewyMm,
          marginRightMm: strona.marginesPrawyMm,
          gutterMm: strona.marginesOprawyMm,
          mirrorMargins: strona.marginesyOdbicia,
          // Nośnik własny niesie też wymiary — jego oznaczenie jest opisem, nie nazwą z wykazu rdzenia.
          ...(strona.nosnik.rodzaj === 'wlasny'
            ? { widthMm: strona.nosnik.szerokoscMm, heightMm: strona.nosnik.wysokoscMm }
            : {}),
          ...autorOperatora,
        });
        if (!wynik.udany || wynik.wynik === undefined) {
          ogloszOdmowe('Nastawy strony z powierzchni', wynik.blad, zaplecze.naZdanie);
          return;
        }
        panelStrony.pokazNastawy(wynik.wynik.pageSetup);
        ogloszSkutek(
          'Nastawy strony z linijki zapisane w rdzeniu — przeżyją zamknięcie okna',
          wynik.wynik.balance,
          wynik.wynik.form,
          zaplecze.naZdanie,
        );
      })();
    },

    zglosWciecia(wciecia) {
      const zadanie: Record<string, unknown> = {};
      if (wciecia.pierwszyWierszMm !== undefined) {
        zadanie['firstLineIndentMm'] = wciecia.pierwszyWierszMm;
      }
      if (wciecia.leweMm !== undefined) zadanie['indentLeftMm'] = wciecia.leweMm;
      if (wciecia.praweMm !== undefined) zadanie['indentRightMm'] = wciecia.praweMm;
      if (Object.keys(zadanie).length === 0) return;
      void nanieStylAkapitu(zadanie, zaplecze.naZdanie);
    },

    zglosTabulator(polozenieMm, zdejmij) {
      const identyfikator = dokument(zaplecze.naZdanie);
      if (identyfikator === null) return;
      void (async () => {
        const wynik = await zaplecze.zrodlo.ustawTabulator({
          documentId: identyfikator,
          positionMm: polozenieMm,
          remove: zdejmij,
          ...zakres(),
          ...autorOperatora,
        });
        if (!wynik.udany || wynik.wynik === undefined) {
          ogloszOdmowe('Tabulator z linijki', wynik.blad, zaplecze.naZdanie);
          return;
        }
        ogloszSkutek(
          `Tabulator na ${polozenieMm} mm ${zdejmij ? 'zdjęty' : 'ustawiony'}; tabulatorów akapitu ` +
            `${wynik.wynik.tabStops.length}`,
          wynik.wynik.balance,
          wynik.wynik.form,
          zaplecze.naZdanie,
        );
      })();
    },

    zglosTabulatory(tabulatory) {
      const identyfikator = dokument(zaplecze.naZdanie);
      if (identyfikator === null) return;
      const zapamietane = tabulatoryAkapitu;
      const zdjete = zapamietane.filter(
        (stary) => !tabulatory.some((nowy) => nowy.milimetry === stary.milimetry),
      );
      const wniesione = tabulatory.filter(
        (nowy) =>
          !zapamietane.some(
            (stary) =>
              stary.milimetry === nowy.milimetry &&
              stary.rodzaj === nowy.rodzaj &&
              stary.znakWiodacy === nowy.znakWiodacy,
          ),
      );
      tabulatoryAkapitu = [...tabulatory];
      if (zdjete.length === 0 && wniesione.length === 0) return;
      void (async () => {
        for (const stary of zdjete) {
          await wyslijTabulator(identyfikator, stary, true);
        }
        for (const nowy of wniesione) {
          await wyslijTabulator(identyfikator, nowy, false);
        }
      })();
    },

    async odczytajPostac() {
      const identyfikator = dokument(zaplecze.naZdanie);
      if (identyfikator === null) return null;
      const wynik = await zaplecze.zrodlo.postac({ documentId: identyfikator });
      if (!wynik.udany || wynik.wynik === undefined) {
        ogloszOdmowe('Odczyt postaci dokumentu', wynik.blad, zaplecze.naZdanie);
        return null;
      }
      zaplecze.naPostac(wynik.wynik.form);
      return wynik.wynik.form;
    },

    async zapiszPostac(postac, tresc, tytul) {
      const identyfikator = dokument(zaplecze.naZdanie);
      if (identyfikator === null) return false;
      const wynik = await zaplecze.zrodlo.zapiszPostac({
        documentId: identyfikator,
        form: postac,
        content: tresc,
        ...(tytul === '' ? {} : { title: tytul }),
        createVersion: true,
        ...autorOperatora,
      });
      if (!wynik.udany || wynik.wynik === undefined) {
        // Nieudany zapis musi być widoczny i nazwany — cichy wskaźnik zapisano byłby błędem najgorszym.
        ogloszOdmowe('Zapis postaci dokumentu', wynik.blad, zaplecze.naZdanie);
        return false;
      }
      ogloszSkutek(
        'Dokument zapisany wraz z POSTACIĄ — arkusz stylów, nastawy strony, sekcje, tabele ' +
          `i aparat przeżyły zapis${wynik.wynik.version === undefined ? '' : `; wersja ${wynik.wynik.version.id}`}`,
        wynik.wynik.balance,
        wynik.wynik.form,
        zaplecze.naZdanie,
      );
      return true;
    },

    async odczytajTekst() {
      const identyfikator = dokument(zaplecze.naZdanie);
      if (identyfikator === null) return null;
      const wynik = await zaplecze.zrodlo.tekst({ documentId: identyfikator, ...zakres() });
      if (!wynik.udany || wynik.wynik === undefined) {
        ogloszOdmowe('Odczyt fragmentu', wynik.blad, zaplecze.naZdanie);
        return null;
      }
      const blokady = wynik.wynik.locks ?? [];
      if (blokady.length > 0) {
        zaplecze.naZdanie(
          `Fragment obejmują blokady: ${blokady
            .map((blokada) => `„${blokada.name}" (znaki ${blokada.rangeStart}–${blokada.rangeEnd})`)
            .join(', ')}. Blokada jest skierowana przeciw modelowi, nie przeciw Operatorowi — ` +
            'chyba że jej zasięg mówi inaczej.',
          true,
        );
      }
      return wynik.wynik.text;
    },

    async zmienTekst(brzmienie, zachowajPostac) {
      const identyfikator = dokument(zaplecze.naZdanie);
      if (identyfikator === null) return false;
      const wybrany = zaplecze.zaznaczenie();
      if (wybrany === null) {
        zaplecze.naZdanie(
          'Zmiana brzmienia dotyczy FRAGMENTU — zaznacz go. Komenda studio.text.edit przyjmuje ' +
            'początek i koniec jako pola obowiązkowe, bo jej sensem jest poprawa bez przepisywania ' +
            'całości.',
          false,
        );
        return false;
      }
      const wynik = await zaplecze.zrodlo.zmienTekst({
        documentId: identyfikator,
        rangeStart: wybrany.poczatek,
        rangeEnd: wybrany.koniec,
        text: brzmienie,
        keepFormat: zachowajPostac,
        ...autorOperatora,
      });
      if (!wynik.udany || wynik.wynik === undefined) {
        ogloszOdmowe('Zmiana brzmienia fragmentu', wynik.blad, zaplecze.naZdanie);
        return false;
      }
      ogloszSkutek(
        `Brzmienie fragmentu (znaki ${wybrany.poczatek}–${wybrany.koniec}) zmienione`,
        wynik.wynik.balance,
        wynik.wynik.form,
        zaplecze.naZdanie,
      );
      return true;
    },

    async odczytajPochodzenia() {
      const identyfikator = dokument(zaplecze.naZdanie);
      if (identyfikator === null) return [];
      const wynik = await zaplecze.zrodlo.pochodzenia({ documentId: identyfikator });
      if (!wynik.udany || wynik.wynik === undefined) {
        ogloszOdmowe('Pochodzenie fragmentów', wynik.blad, zaplecze.naZdanie);
        return [];
      }
      return wynik.wynik.entries;
    },

    magazynWidoku,
    postacZapamietana: () => postacOstatnia,
  };
}

/** Zestawienie postaci dokumentu podaje liczby, nie zawartość: arkusz stylów, sekcje i numer porządkowy postaci potwierdzają, że dokument niesie to, co Operator zapisał. */
function opiszPostac(postac: StudioDocumentForm): string {
  return (
    `stylów ${postac.styles?.length ?? 0} · sekcji ${postac.sections?.length ?? 0} · bloków ` +
    `${postac.blocks?.length ?? 0} · tabel ${postac.tables?.length ?? 0} · obiektów ` +
    `${postac.objects?.length ?? 0} · list ${postac.lists?.length ?? 0} · pozycji aparatu ` +
    `${postac.apparatus?.length ?? 0} · pól ${postac.fields?.length ?? 0} · blokad ` +
    `${postac.locks?.length ?? 0}` +
    `${postac.revision === undefined ? '' : ` · numer porządkowy postaci ${postac.revision}`}`
  );
}
