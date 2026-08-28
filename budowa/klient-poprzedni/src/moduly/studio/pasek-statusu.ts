import { StudioDocumentFormat } from '../../../../shared/contract';
import { opiszLiczniki, policzTresc } from './liczniki-dokumentu';
import type { StanStudio } from './stan-studio';

/** Zwraca nazwę formatu dokumentu widoczną dla operatora; format spoza kontraktu zostaje w postaci surowej. */
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

      // Brak dokumentu i format nieznany to dwie różne rzeczy: pierwsza to brak treści, druga brak danych.
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
