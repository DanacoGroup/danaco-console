import { pokazKomunikat } from '../../aplikacja/komunikaty';

/**
 * Funkcje wykazu okien modułu Assistant, dla których kontrakt nie ma drogi.
 * Jedyną odpowiedzialnością tego pliku jest nazwanie braku, a każde zdanie
 * wymienia komendę albo pole, którego brakuje.
 */
export const BRAKI = {
  szukanieZnaczeniem:
    'Wyszukiwanie w dzienniku idzie po treści wpisu, nie po znaczeniu. ' +
    'knowledge.search umie szukać po sensie, ale jego zakresy (KnowledgeScope) to ' +
    'biblioteka, historia rozmów i pliki przestrzeni roboczej — dziennik asystenta nie ' +
    'jest żadnym z nich, a assistant.activity.list nie przyjmuje frazy. ' +
    'Drogę otworzy zakres dziennika w KnowledgeScope albo pole zapytania w komendzie ' +
    'dziennika.',
  bibliotekaSkrotow:
    'Siatka szybkich akcji jest katalogiem rdzenia (action.list). Kontrakt nie ' +
    'ma komendy zapisu własnej biblioteki poleceń szybkich, więc zestawu nie da się ' +
    'dziś zmienić z okna.',
  brakOkna:
    'Rdzeń nie oddał okna modułu Assistant dla tej sesji (window.list). ' +
    'assistant.voice.command wymaga pola windowId, więc polecenia nie ma dokąd wysłać.',
  bramaPotwierdzen:
    'Bramy potwierdzeń akcji o skutkach ubocznych (wysyłka, płatność, usunięcie) kontrakt ' +
    'nie ma czym poprowadzić: AssistantAction nie niesie ani znacznika nieodwracalności, ' +
    'ani stanu „oczekuje na decyzję", a AssistantActionControl zna pięć wartości — none, ' +
    'pause, resume, cancel, retry — i żadna z nich nie znaczy „zatwierdzam". Blokada ' +
    'zbudowana wyłącznie w oknie zatrzymywałaby zlecenie u siebie, a rdzeń prowadziłby ' +
    'je dalej. Wstrzymanie zlecenia działa i jest tym, czym dziś dysponuje Operator.',
  brakSesji:
    'Powłoka nie podała karty sesji, dla której moduł został wczytany. ' +
    'memory.list, memory.toggle i session.tool.* wymagają pola sessionId, więc nie ma ' +
    'czego czytać ani gdzie zapisać.',
  wgranieDokumentu:
    'Dokumenty wchodzą do platformy przez moduł Library, nie przez to okno. ' +
    'knowledge.index obejmuje wskaźnikiem znaczenia to, co w bibliotece już leży, ' +
    'a drugie wejście dla plików znaczyłoby dwa repozytoria treści Operatora. ' +
    'Wgraj plik w module Library, po czym zbuduj tu wskaźnik.',
  zakresNarzedzia:
    'Zakresu uprawnień i limitu wywołań narzędzia kontrakt nie ma dla profilu asystenta. ' +
    'ToolCatalogEntry mówi, czym pozycja jest i czy da się ją dołożyć (attachable, ' +
    'attached), a session.tool.attach doklada ją do karty sesji w całości — bez granic ' +
    'argumentów i bez licznika wywołań. Zakres per ekspert prowadzi osobna rodzina ' +
    'agent.connector.* modułu Agents i nie sięga ona profilu asystenta.',
  makroWlasne:
    'Makra asystenta kontrakt nie przechowuje: nie ma ani komendy zapisu makra, ani ' +
    'wykazu poleceń szybkich profilu. Jedynym trwałym miejscem dla sekwencji kroków jest ' +
    'automatyka modułu Automations (automation.workflow.save), więc makro zapisuje się ' +
    'właśnie tam — przyciskiem „→ Wyślij do Automations".',
  makroYaml:
    'Edytor przyjmuje wyłącznie JSON. YAML wymagałby biblioteki, której klient nie ma, ' +
    'a dołożenie zależności jest decyzją architektoniczną Właściciela, nie okna. ' +
    'Kroki w JSON idą do rdzenia dokładnie tak, jak wymaga tego pole steps komendy ' +
    'automation.workflow.save.',
} as const;

/**
 * Odpowiedź na naciśnięcie kontrolki, dla której kontrakt drogi nie ma: okno
 * pokazuje komunikat wagi ostrzegawczej ze zdaniem nazywającym brak, zamiast
 * milczeć albo wygaszać kontrolkę.
 */
export function zglosBrak(tytul: string, tresc: string): void {
  pokazKomunikat({ tytul, tresc, waga: 'ostrz' });
}
