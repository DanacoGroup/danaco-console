import { StudioDocumentFormat } from '../../../../shared/contract';
import { opiszLiczniki, policzTresc } from './liczniki-dokumentu';
import type { StanStudio } from './stan-studio';

/**
 * Pasek statusu dokumentu — element warstwy pierwszej Studio Editora.
 *
 * Opracowanie stawia go pod obszarem treści i żąda od niego trzech rzeczy:
 * stanu zapisu, liczby słów i numeru wersji. Wszystkie trzy da się powiedzieć
 * prawdziwie z tego, co moduł już ma, i żadna nie wymaga wywołania rdzenia.
 *
 * Stan zapisu poznaje się po różnicy między treścią roboczą a zaakceptowaną:
 * bufor edytora i treść ustalona to w tym module dwa osobne pola
 * (`pola-stanu.ts`), więc „niezapisane zmiany" jest tu faktem odczytanym,
 * a nie znacznikiem ustawianym ręcznie przy każdym naciśnięciu klawisza.
 *
 * Numer wersji bierze się z `StudioDocument.versionId`. Dokument przed
 * pierwszym zapisem wersji nie ma i pasek mówi to wprost, zamiast pokazywać
 * zero — zero wyglądałoby na wersję o numerze zero.
 *
 * Format dokumentu stoi tu od scalenia okien: nosił go wskaźnik paska narzędzi
 * edytora, a pasek narzędzi rozszedł się na wstążkę okna pracy. Format jest
 * cechą dokumentu, nie czynnością, więc jego miejsce jest w pasku stanu —
 * `StudioDocument.format` przychodzi w odpowiedzi `studio.document.open`.
 */

/** Nazwa formatu dla Operatora; format spoza kontraktu zostaje w postaci surowej. */
function nazwaFormatu(format: string): string {
  const nazwy: Record<string, string> = {
    [StudioDocumentFormat.Pdf]: 'PDF',
    [StudioDocumentFormat.Docx]: 'DOCX',
    [StudioDocumentFormat.Txt]: 'TXT',
    [StudioDocumentFormat.Markdown]: 'Markdown',
  };
  return nazwy[format] ?? format;
}
export interface PasekStatusu {
  element: HTMLElement;
  odswiez(): void;
}

export function utworzPasekStatusu(stan: StanStudio): PasekStatusu {
  const zapis = document.createElement('span');
  zapis.className = 'dn-plakietka ms-status__zapis';

  const liczniki = document.createElement('span');
  liczniki.className = 'ms-status__liczniki';

  const wersja = document.createElement('span');
  wersja.className = 'dn-plakietka ms-status__wersja';

  const format = document.createElement('span');
  format.className = 'dn-plakietka ms-status__format';

  const element = document.createElement('p');
  element.className = 'ms-status';
  element.setAttribute('aria-label', 'Pasek statusu dokumentu');
  element.append(zapis, liczniki, wersja, format);

  return {
    element,

    odswiez() {
      const dokument = stan.dokument();
      const robocza = stan.trescRobocza();
      const rozbieznosc = robocza !== stan.trescZaakceptowana();

      // Stan nigdy samym kolorem: obok odmiany plakietki idzie zawsze napis.
      element.dataset['zapis'] = rozbieznosc ? 'niezapisane' : 'zapisane';
      zapis.className = rozbieznosc
        ? 'dn-plakietka dn-plakietka--ostrzezenie ms-status__zapis'
        : 'dn-plakietka dn-plakietka--sukces ms-status__zapis';
      zapis.textContent = rozbieznosc ? 'niezapisane zmiany' : 'zapisane';
      zapis.title = rozbieznosc
        ? 'Bufor edytora różni się od treści ustalonej. Zapis w Studio Editorze utrwali różnicę i założy wersję.'
        : 'Bufor edytora zgadza się z treścią ustaloną przez rdzeń.';

      liczniki.textContent = opiszLiczniki(policzTresc(robocza));

      // Brak dokumentu i format nieznany to dwie różne rzeczy: pierwsza znaczy
      // „nie ma czego formatować", druga — „rdzeń formatu nie podał".
      format.textContent =
        dokument === null
          ? 'format: brak dokumentu'
          : `format z rdzenia: ${nazwaFormatu(dokument.format)}`;

      const numerWersji = dokument?.versionId ?? '';
      wersja.textContent =
        dokument === null
          ? 'bez dokumentu'
          : numerWersji === ''
            ? 'dokument jeszcze bez wersji'
            : `wersja bieżąca: ${numerWersji}`;
    },
  };
}
