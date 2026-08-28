# Uzasadnienia pakietów

| | |
|---|---|
| **Produkt** | Danaco Console |
| **Przeznaczenie** | Uzasadnienia projektowe i zastrzeżenia przeniesione z nagłówków kodu, których pełna treść przekraczała dopuszczalną długość komentarza |

Dokument gromadzi rozdziały nazwane ścieżką pliku źródłowego. Każdy rozdział wylicza identyfikatory (typ, pole, stała, funkcja) wraz z uzasadnieniem, które nie mieściło się w nagłówku komentarza tego pliku.

---

### budowa/shared/contract.go

Uzasadnienia i zastrzeżenia projektowe przeniesione z komentarzy kontraktu komunikacji WebSocket, którym nagłówek kodu nie miał miejsca.

**MessageType** — Jest aliasem napisu, wiec stale kontraktu wchodza wprost tam, gdzie warstwa protokolu operuje na string.

**ExecutionEnv** — Nazwa konkretnego hosta nie jest wartoscia wyliczenia — idzie ustawieniem na poziomie zasiegu window

**PermissionMode** — Wartosci odpowiadaja doslownie przelacznikowi --permission-mode kanalu glownego; zweryfikowane z claude --help 2026-08-11.

**ChunkKind** — Ten sam slownik opisuje rodzaj tresci wiadomosci zapisanej w bazie

**ChunkKind.ChunkKindFinal** — ZASTEPUJE tresc zlozona z dotychczasowych fragmentow tekstowych, nie doklada sie do niej. Niesie ja koperta domykajaca (done=true); odbiorca podmienia tekst wpisu na tresc tego fragmentu. Wartosc PRZELOTOWA — tresc ostateczna utrwala sie jako `content` wiadomosci o rodzaju `tekst`, wiec sam ten rodzaj do kolumny nigdy nie trafia

**ConfigScope.ConfigScopeApplication** — Poziom NAJSZERSZY: przegrywa z kazdym wezszym zapisem, a zapisu na nim nie ma czym zawezic, bo bytu nie ma (scopeId pusty). Tedy przelacznik ustalenie budowy dochodzi do warstwy nasluchu bez restartu rdzenia

**NavigationKind** — Trzy srodowiska wystawiaja liste modulow, MultitaskingAI wystawia panel orkiestracji i nie udostepnia modulow

**ConfigAxis** — Os jest prostopadla do poziomu zasiegu: poziom mowi jak wasko (okno -> ... -> globalny), os mowi dla czego (platforma, model, konto). Klucz rozstrzygania jest zlozony: klucz + poziom + byt poziomu + os + byt osi. Brak osi znaczy platform

**ConfigAxis.ConfigAxisPlatform** — Wartosc domyslna osi — os pominieta znaczy platform

**SettingValueType** — Klient buduje z niego kontrolke formularza i nie zna zadnego klucza z osobna

**AccessPointKind** — Punkt dostepu nie jest srodowiskiem: srodowisko pozostaje profilem widocznosci modulow w bocznej nawigacji

**IdentityLayer** — Wartosci odpowiadaja stalym silnika nakladki

**IdentityMode** — Tozsamosc jest ZAMIENIANA, nie dolaczana — wartoscia domyslna klucza tozsamosc.tryb_domyslny jest ZASTAP. Wartosci sa sterowaniem istniejacego silnika nakladki, nie jego powtorzeniem

**IdentityMode.IdentityModeZASTAP** — Tryb zastapienia nie podaje --append-system-prompt

**LoopStopReason** — Zatrzymanie nigdy nie jest ciche; wartosci odpowiadaja stalym session.PowodZatrzymania

**LoopStopReason.LoopStopReasonCompleted** — Jedyna wartosc mowiaca skonczone, a nie przerwane — po niej bieg podejmuje sie sam przy kolejnym koncu tury wykonawcy, wiec ukonczenie nie jest brama akceptacji; odpowiada session.ZatrzymanieUkonczenie

**SessionConfigArea** — Obszar jest jednostka zapisu: komenda config.session.set bierze wykaz obszarow, nie wykaz pol. Warstwy tozsamosci, profilu i ekspertyzy nie sa obszarami — sa polami obszaru systemPrompt

**ProviderTransport** — Wartosc wyznacza adapter, ktory tlumaczy konfiguracje sesji na powierzchnie dostawcy — sam model konfiguracji pozostaje ten sam

**AdapterSurface** — Interfejs nie jest budowany wokol nazw powierzchni — sluza wylacznie diagnostyce i deklaracji zdolnosci

**SessionConfigSource** — Konfiguracja sesji SPINA istniejace rodziny, nie zastepuje ich — pole mowi, ktory rejestr byl zrodlem prawdy

**WorkingDirectoryDegradation** — Rozejscie ma byc widoczne w oknie, nie milczace

**DesignAssetKind.DesignAssetKindDocument** — Wartosc dopisana 15.08.2026 — narzedzia dokumentowe wytwarzaly pliki, ktore trzeba bylo opisywac jako wektor, bo enum nie znal dokumentu

**DesignAssetKind.DesignAssetKindAudio** — Wartosc dopisana 15.08.2026 — narzedzia mediow wytwarzaly dzwiek, ktory trzeba bylo opisywac jako obraz

**DesignAssetKind.DesignAssetKindVideo** — Wartosc dopisana 15.08.2026 z tego samego powodu co audio

**AgentPermissionGroup.AgentPermissionGroupModules** — Zakres szczegolowy wpisu niesie kod modulu z rodziny `module.*`, wiec wykaz modulow nie jest tu powielany

**MemoryLevel** — WYLACZENIE PAMIECI NIE JEST PIATA WARTOSCIA TEGO WYLICZENIA — jest PUSTYM ZBIOREM poziomow w polu `Agent.memoryLevels`. Piata wartosc `disabled` obok czterech poziomow pozwalalaby zapisac stan sprzeczny (`[session, disabled]`), ktorego rdzen musialby rozstrzygac zgadywaniem; zbior pusty jest jedynym ksztaltem, w ktorym "wylaczona" nie da sie postawic obok "wlaczona". Wartosc wyjsciowa definicji eksperta jest nadpisywana przy przypisaniu do projektu i przy przypisaniu do roli — pierwszenstwo ma zasieg najbardziej szczegolowy

**AgentVisibility** — Ekspert `project` NIE MA wlasnego wskazania projektu: widoczny jest w tych projektach, do ktorych zostal przypisany komenda `workspace.agent.assign` — drugie miejsce zapisu przynaleznosci byloby druga prawda o tym samym

**ExtensionOrigin** — RDZEN NIE ROZROZNIA ZRODLA: ladowanie, kontrakt integracji i uzycie operacyjne sa dla obu zrodel te same. Jedyny wyjatek zapisany przez ustalenie budowy to STAN WYJSCIOWY PRZY REJESTRACJI — `danaco` staje wlaczone, `personal` wylaczone. Poza tym jednym miejscem zrodlo jest faktem do pokazania Operatorowi, nie przeslanka rozgalezienia zachowania

**AuthMethodKind** — Trzy, wzorem Danaco HUB: haslo Operatora, PIN i Windows Hello. Bez OAuth, bez federacji, bez rotacji kluczy, bez MFA. Konta uzytkownika NIE MA — jedyny Operator jest bezimienny

**AuthMethodKind.AuthMethodKindHello** — NIEDOSTEPNA POD BIEZACYM POCHODZENIEM: WebAuthn wywodzi rp_id z pochodzenia dokumentu, a powloka natywna podaje interfejs z http://127.0.0.1:<port>, czyli z adresu IP, ktory domena nie jest. Metoda staje sie dostepna dopiero, gdy interfejs jest podawany z pochodzenia z domena po https

**AdvisorSelection** — Dwie wartosci, bo drogi sa dwie: wskazanie Operatora albo sufit sily. Trzeciej — eskalacji z inicjatywy modelu — nie ma i nie bedzie: rozstrzygniecie ustalenie budowy z 14.08.2026 mowi, ze bez wskazania Operatora model kontraktuje agenta o parametrach takich samych albo nizszych

**AdvisorSelection.AdvisorSelectionOperatorIndication** — Wskazanie znosi sufit sily, bo pochodzi od czlowieka, nie od modelu

**AdvisorSelection.AdvisorSelectionStrengthCeiling** — Tak konczy sie kazda konsultacja, o ktora poprosil sam model

**ActorKind** — Bez tego rozroznienia Operator widzi skutek, ale nie widzi REKI — a asystent główny prace za niego staje sie nieodrozniahy od niego samego

**SlashEntryKind** — To NIE SA dwa mechanizmy: jedno pole, jeden wykaz, jedno filtrowanie — rozne sa skutki wyboru. Klasyfikacja musi istniec w danych, bo pozycja rodzaju `tool` POSZERZA zestaw narzedzi sesji, a pozycja rodzaju `action` wykonuje czynnosc aplikacji i zestawu nie dotyka

**SessionToolSource** — Komenda po ukosniku jest JEDYNA droga poszerzenia zestawu w trakcie pracy (rozdz. 2 „Wylacznosc"), a to pole mowi, kto te droge przeszedl — Operator z klawiatury czy asystent dzialajacy za niego. Nie jest to druga droga dolozenia, tylko zapis sprawcy

**AgentPermissionAction** — Dzisiejszy `AgentPermission` niesie grupe, zakres i wartosc logiczna, wiec zawezenie „konektor tylko do odczytu" nie ma jak zostac zapisane — wartosc logiczna rozstrzyga wylacznie „wolno albo nie wolno"

**AgentModuleScope** — Wartosci odpowiadaja przykladom wymienionym w opracowaniu rozdz. 9.2): zapis plikow w module Library, wykonywanie polecen w module Terminal

**TelemetryFormat** — Rozlaczny z ExportFormat, ktory sluzy dokumentom Operatora (pdf, docx, markdown) i formatow maszynowych nie niesie

**UsageDimension** — Jeden wymiar na zadanie: zestawienie po dwoch wymiarach naraz jest dwoma zadaniami, nie jednym

**AlertMetric** — Kazda wartosc ma zrodlo w danych, ktore platforma juz przechowuje — reguly nad miara bez zrodla nie da sie ewaluowac

**AlertChannel** — Poza droga w aplikacji kazda wymaga skonfigurowanej integracji; brak integracji nie wstrzymuje wyzwolenia, tylko dostarczenie

**LibraryFileStatus** — Archiwizacja jest odwracalna: zasob zarchiwizowany znika z wykazu domyslnego, nie z repozytorium

**LibraryRuleKind** — Obie skladaja plik z kolekcja, ale z innej strony: regula kolekcji liczy zawartosc z warunku, regula naplywu wciaga zasob z zewnatrz

**LibraryPreservationKind** — Kazda jest odrebnym standardem branzowym, nie wariantem jednego zapisu

**LibraryDuplicateKind** — Bez niej wynik dokladny i wynik przyblizony wygladaja tak samo

**LibraryFieldKind** — Rozstrzyga postac kontrolki i sprawdzian wartosci, nie blokade zapisu

**LibraryRetentionAction** — Twarde usuniecie nie zachodzi samoczynnie

**WorkspaceTaskStatus** — Kolumny tablicy sa nastawa Operatora (rozdz. 9 opracowania) i odwzorowuja sie na te stany; kolumna wlasna nie zaklada stanu nowego

**WorkspaceDependencyKind** — Dwie wartosci, bo tylko dwie ma pokrycie w opracowaniu: nastepstwo poprzednik-nastepnik oraz blokada

**WorkspaceEntityKind** — Jedno wyliczenie dla wyszukiwania, komentarzy, wezlow grafu i osi czasu — te same byty pod jedna nazwa, bez czterech list synonimow

**StudioActionKind.StudioActionKindFieldChange** — Osobno od objectChange, bo pole odlozone jako zmiana obiektu nie da sie cofnac samo — Operator cofalby wtedy cala zmiane obiektu razem z nia

**AodEventClass** — Klasa jest jednoczesnie zakresem wyciszenia (rozdz. 3.5, wiersz trzeci): wyciszona klasa nie tworzy sugestii do chwili zniesienia, a pozostale dzialaja bez zmian

**AodMuteKind** — Tryb cichy NIE JEST czwarta wartoscia tego wyliczenia: jest trybem obecnosci (rozdz. 9.2) i jedzie ustawieniem konfiguracji, nie wpisem wyciszenia. Wyjatek wagi krytycznej takze nie jest wyciszeniem — jest regula przebijajaca kazde z trzech

**AodMuteKind.AodMuteKindTimed** — Chwile konca niesie pole endsAt

**AodMuteScope** — Wyciszenie czasowe obowiazuje cala platforme i zakresu nie niesie (pole scope puste); wyciszenie kontekstowe wskazuje modul albo karte sesji, wyciszenie klasy zdarzen — klase

**MessageType.CommandSpeechAvailabilityGet** — Brak silnika to ODPOWIEDZ, nie awaria — klient nie pokazuje mikrofonu, ktory nic nie nagra

**MessageType.CommandSpeechTranscribe** — Dzwiek NIE OPUSZCZA maszyny, na ktorej stoi silnik

**MessageType.CommandSessionFocus** — Asystent przestawia ognisko OPERATOROWI polem targetClientId — inaczej jego posuniecia nie byly widoczne na ekranie

**MessageType.CommandDesignAssetGenerate** — Tresc trafia do magazynu rdzenia pod suma kontrolna, wiec zasob nie zalezy od zadnego pliku zewnetrznego. Brak kanalu obrazowego albo brak poswiadczenia to ODMOWA NAZYWAJACA BRAK — rdzen nigdy nie zaklada zasobu bez bajtow obrazu

**MessageType.CommandAgentLayerSet** — Warstwa DOPISUJE sie do promptu systemowego, nie zastepuje go — trybu podania sie tu nie wybiera

**MessageType.CommandMemoryDetach** — Odpiecie NIE JEST wylaczeniem: zweza zasieg samego wpisu, a nie wstrzymuje go w cudzym projekcie ani module — do tego sluzy memory.disable.set

**MessageType.CommandComponentAssign** — UWAGA: znaczenie przypisania nie zostalo rozstrzygniete przez ustalenie budowy

**MessageType.CommandExtensionInstall** — UWAGA: znaczenie instalacji — pobranie paczki, zarejestrowanie adresu czy zapis punktu dostepu — nie zostalo rozstrzygniete przez ustalenie budowy

**MessageType.CommandExtensionUninstall** — UWAGA: znaczenie odinstalowania jest zwiazane z nierozstrzygnietym znaczeniem instalacji

**MessageType.CommandAdvisorConsult** — Rada NIE JEST WIAZACA i nie wykonuje pracy za pytajacego — odpowiedzialnosc za wynik zostaje przy agencie, ktory pyta. Kim jest doradca, rozstrzyga RDZEN z danych okna, nie zadanie: pole requestedAdvisor jest PROSBA i podlega sufitowi sily. Prosba o model silniejszy od kanalu pytajacego bez wczesniejszego wskazania Operatora jest ODMAWIANA, a nie zamieniana po cichu na slabszego doradce — model z wlasnej inicjatywy po model silniejszy nie siega (rozstrzygniecie ustalenie budowy z 14.08.2026). Kazda odbyta konsultacja rozglasza sie zdarzeniem advisor.consulted

**MessageType.CommandConfigExplainGet** — Transparentnosc, nie bramka

**MessageType.CommandConfigWindowOpen** — UWAGA: wykaz okien obiecuje te komende, nie mowiac, co rdzen ma przy niej robic poza podaniem zakresu

**MessageType.CommandTerminalOutputRead** — Komenda `terminal.command.exec` konczy sie w chwili STARTU procesu (kompilacja trwa dluzej niz kazde sensowne oczekiwanie na odpowiedz, a rozlaczenie klienta nie ma prawa jej przerwac), wiec wyniku niesc nie moze i nigdy nie bedzie mogla. Tedy model czyta, co polecenie wypisalo. Zrodlem jest ten sam dziennik zbiorczego wyjscia, na ktorym stoi `terminal.output.stream` — drugiej pompy nie ma. Historia zyje jeden bieg rdzenia: po restarcie wyjscie jest puste i odpowiedz mowi to wprost

**MessageType.CommandModelChannelSet** — Kanaly zaklada i zmienia rodzina channel.*; ta komenda wybiera jeden z zalozonych

**MessageType.CommandAodSuggestion** — UWAGA: wykaz okien obiecuje te pozycje slowem podpowiedz, nie rozstrzygajac, czy jest to odczyt Operatora, czy zdarzenie wypychane przez rdzen

**MessageType.CommandAuthRegister** — Konto powstaje w stanie NIEPOTWIERDZONYM i pozostaje w nim do chwili potwierdzenia adresu komenda auth.verify — dopiero potwierdzenie wydaje urzadzeniu token dostepu. Sesji ta komenda NIE zaklada. Wykonalna tylko raz; potem odmawia trwale. Powtorzenie hasla jest sprawa formularza klienta, nie kontraktu

**MessageType.CommandAuthVerify** — Przenosi konto ze stanu niepotwierdzonego do potwierdzonego, po czym wydaje urzadzeniu token dostepu — to jest moment, w ktorym Operator wchodzi do platformy po raz pierwszy. Droga potwierdzenia przychodzi listem na adres podany przy rejestracji i wygasa; wygasla droga odmawia i pozwala poprosic o nowa

**MessageType.CommandAuthRecover** — Serwer wysyla na wskazany adres droge potwierdzenia tozsamosci; nowe haslo ustawia sie komenda auth.reset. Odpowiedz nie zdradza, czy adres pasuje do konta — inaczej komenda mowilaby obcemu, jaki adres ma Operator

**MessageType.CommandAuthReset** — Zastepuje dotychczasowy skrot hasla przechowywany poza baza danych i uniewaznia tokeny dostepu wydane przed odzyskaniem — kazde powiazane urzadzenie loguje sie ponownie, a zmiana idzie w swiat zdarzeniem device.changed. Nowej encji Konto nie tworzy: zmienia sie wylacznie material uwierzytelniajacy

**MessageType.CommandAuthLogin** — Metodami dzialajacymi dzis sa HASLO i PIN; hello czeka na pochodzenie z domena po https. Login niesie sie przy metodzie password, bo konto ma nazwe ustalenie budowy nadana przy rejestracji; metody wlasciwe urzadzeniu tozsamosci nie potrzebuja, bo wskazuje ja material na urzadzeniu. Odmowa wraca kodem bledu not_authenticated. Proba nieudana NIE odmawia nastepnej — nakłada na nia ZWLOKE rosnaca wykladniczo (250 ms, 500, 1000, 2000, 4000, dalej rowno 5000 ms), zerowana przy pierwszym udanym wejsciu. Operator, ktory pomylil haslo trzy razy, wchodzi za czwartym; tylko czeka. Zadnego progu prob i zadnej odmowy 'za duzo prob' tu nie ma i nie bedzie — to byloby bramkowanie (rozstrzygniecie ustalenie budowy z 14.08.2026)

**MessageType.CommandAuthMethodAdd** — Czynnosc USTAWIEN, wykonalna po zalogowaniu; hasla ta komenda nie zaklada, bo kotwica powstaje przy auth.register. Klucz Hello nie opuszcza TPM urzadzenia, wiec na kazdej maszynie zaklada sie go osobno — i tylko wtedy, gdy interfejs jest podawany z pochodzenia z domena po https. Dzis zalozyc mozna PIN

**MessageType.CommandAuthMethodRemove** — Ani ostatniej metody, ani hasla zdjac sie nie da — haslo jest kotwica bramki

**MessageType.CommandAuthTokenRefresh** — W Danaco HUB dzieje sie to samo przy kazdym zadaniu HTTP; tu jest komenda, bo transportem jest jedno gniazdo WebSocket

**MessageType.CommandDeviceList** — Konto jest jedno, urzadzen dowolnie wiele: komputery, telefony, tablety. Kazde niesie wlasny token dostepu, wiec wykaz jest miejscem, w ktorym Operator widzi, co ma dostep do platformy

**MessageType.CommandDeviceRevoke** — Zmiana idzie do wszystkich polaczonych urzadzen zdarzeniem device.changed, wiec Operator widzi skutek natychmiast na pozostalych ekranach. Tej samej drogi uzywa odzyskanie konta, ktore uniewaznia tokeny wydane przed zmiana hasla

**MessageType.CommandAppsDeploymentList** — Bez tej komendy historia wdrozen znikala po odswiezeniu okna, bo klient odbudowywal ja wylacznie ze zdarzen biezacej sesji

**MessageType.CommandAppsArchitectureGet** — Odpowiednik odczytu dla apps.architecture.define — bez niego okno po odswiezeniu nie wie, co Operator wczesniej zdefiniowal

**MessageType.CommandAppsWorkspaceList** — Odpowiednik odczytu dla apps.workspace.update

**MessageType.CommandDesignAssetUpload** — Tresc wchodzi do magazynu rdzenia pod suma kontrolna — zasob przestaje zalezec od pliku, ktory Operator moze nadpisac albo skasowac. Uzyj, gdy zasob JUZ ISTNIEJE: Operator ma plik albo wskazuje odsylacz. Zasob majacy dopiero powstac zamawia sie komenda design.asset.generate

**MessageType.CommandDesignAssetRemove** — Rodzaj zmiany deleted zdarzenia design.asset.changed nie mial dotad zadnego nadawcy

**MessageType.CommandImageTransform** — Wynikiem jest NOWY zasob w magazynie — zrodlo zostaje nietkniete, wiec model moze probowac bez ryzyka

**MessageType.CommandImageAdjust** — To jest RETUSZ, o ktory model prosi w rozmowie

**MessageType.CommandImageConvert** — Sluzy przygotowaniu zasobu do wydania — model uzywa jej przed osadzeniem grafiki na stronie

**MessageType.CommandMediaTranscode** — Wynikiem jest nowy zasob w magazynie

**MessageType.CommandDocumentConvert** — Model uzywa jej, gdy Operator prosi o dokument, a nie o tekst w oknie

**MessageType.CommandDocumentTextExtract** — Model uzywa jej, zeby PRZECZYTAC to, co dostal jako plik

**MessageType.CommandArchivePack** — Sluzy wydaniu pracy Operatorowi jednym plikiem

**MessageType.CommandMailAccountList** — Bez skonfigurowanej skrzynki rodzina mail.* odmawia, nazywajac brak — rdzen poczty nie zmysla

**MessageType.CommandMailMessageList** — Uzyj do odnalezienia sprawy, o ktorej mowi Operator —na przykład listu od konkretnej osoby z dzisiaj

**MessageType.CommandMailMessageGet** — Uzyj, gdy masz przeanalizowac list, a nie tylko go odnalezc. Zalaczniki trafiaja do magazynu rdzenia i wracaja jako zasoby, wiec model moze je dalej przetworzyc narzedziami obrazu i dokumentow

**MessageType.CommandMailDraftSave** — Szkic jest krokiem POSREDNIM: Operator widzi go w swojej poczcie, zanim cokolwiek wyjdzie w swiat

**MessageType.CommandMailSend** — TO JEST JEDYNA KOMENDA RODZINY, KTOREJ SKUTEK WYCHODZI POZA MASZYNE OPERATORA i nie da sie go cofnac

**MessageType.CommandMailAccountAdd** — Rdzen nie zaklada konta pocztowego i nie stawia serwera — bierze skrzynke, ktora juz istnieje: na urzadzeniu albo w chmurze. Poswiadczenie idzie do sejfu, nie do bazy

**MessageType.CommandMailAccountDiscover** — Niczego nie podpina i nie siega po haslo — oddaje to, co znalazl, zeby Operator nie przepisywal nastaw recznie

**MessageType.CommandMailAccountRemove** — Skrzynki u dostawcy nie tyka

**MessageType.CommandImageUpscale** — Zastepuje wyspecjalizowane modele powiekszajace; bez zainstalowanego silnika odmawia, nazywajac brak — nigdy nie oddaje zwyklego rozciagniecia jako powiekszenia

**MessageType.CommandImageBackgroundRemove** — Zastepuje wyspecjalizowane narzedzia wycinania; bez silnika odmawia, nazywajac brak

**MessageType.CommandKnowledgeIndex** — Wyszukiwanie po slowach juz dziala (library.file.search); to jest droga do wyszukiwania po SENSIE

**MessageType.CommandKnowledgeImageSearch** — Rodzina knowledge.* prowadzi dotad wylacznie tekst; ta komenda jest jej osia obrazu i przeglada obrazy biblioteki Operatora

**MessageType.CommandSessionToolAttach** — Definicji eksperta NIE RUSZA (rozstrzygniecie ustalenie budowy, rozdz. 5): dolozenie zyje w stanie sesji, przezywa rozlaczenie klienta i konczy sie wraz z sesja albo z usunieciem rozmowy

**MessageType.CommandSessionToolDetach** — Zestaw wraca do podstawy z definicji eksperta; sama definicja nie zmienia sie ani o joto, bo nigdy nie byla zmieniana

**MessageType.CommandSessionToolList** — Zestaw narzedzi tury to definicja eksperta PLUS te dolozenia — bez tego odczytu druga polowa zestawu bylaby niewidoczna

**MessageType.CommandToolsCatalogList** — Wykaz liczy setki pozycji, wiec zadanie niesie zawezenie tekstem, rodzajem i grupa; dopasowanie idzie takze SRODKIEM nazwy, bo przy przedrostkach zrodla szukanie od poczatku byloby bezuzyteczne (rozdz. 7.2)

**MessageType.CommandAgentConnectorList** — Odpowiednik `agent.plugin.list` dla drugiego z dwoch bytow

**MessageType.CommandAgentVersionGet** — `AgentVersion` niesie etykiete, opis zmiany, sprawce i sume kontrolna, ale NIE niesie tresci — bez tej komendy podgladu wersji nie ma z czego zlozyc, a przywrocenie jest jedynym sposobem zobaczenia, co w wersji stalo

**MessageType.CommandAgentAssignmentList** — `workspace.agent.assign` zapisuje przynaleznosc, ale zadna komenda nie odczytuje jej OD STRONY EKSPERTA, wiec biblioteka nie ma skad wziac liczby

**MessageType.CommandAgentModulesSet** — LISTA PUSTA `[]` ZNACZY BRAK OGRANICZENIA, czyli dostepnosc wszedzie, i jest to stan wyjsciowy; pominiecie pola znaczy „bez zmiany". Ta sama zasada co przy `memoryLevels`, gdzie zbior pusty tez niesie znaczenie wlasne

**MessageType.CommandAgentIsolationSet** — Wartosc obowiazuje wszedzie, gdzie ekspert dziala, dopoki nie nadpisze jej regula zapisana na poziomie zasiegu bardziej szczegolowym

**MessageType.CommandAgentPolicyGet** — Transparentnosc, nie bramka — komenda niczego nie zapisuje i niczego nie rozstrzyga

**MessageType.CommandAgentSubagentSet** — Granica pietnastu jest parametrem technicznym platformy, nie decyzja produktowa — ta sama, ktora zna `subagent.spawn`

**MessageType.CommandAgentPermissionRemove** — `agent.permission.set` wylacznie USTAWIA wartosc wpisu, wiec po pierwszym zawezeniu wiersz zakresu zostaje w wykazie na zawsze, choc jego wartosc wrocila do przyznanej

**MessageType.CommandChannelCheck** — Narzedzie pomocnicze, nie bramka: wynik nie warunkuje zapisu eksperta ani wyslania tury

**MessageType.CommandChannelCredentialStatus** — NIE ODDAJE TRESCI POSWIADCZENIA I ODDAWAC JEJ NIE MOZE

**MessageType.CommandExtensionSearch** — Dopelnia extension.list, ktory zawezal wylacznie rodzajem i stanem zainstalowania; dopasowanie idzie takze SRODKIEM nazwy, bo pozycje niosa przedrostki zrodla

**MessageType.CommandExtensionDetailGet** — Extension niesie sam naglowek pozycji, wiec karta szczegolow nie ma dzis z czego powstac

**MessageType.CommandExtensionCollectionApply** — Wynik jest bilansem, nie potwierdzeniem: pozycja odrzucona wraca z powodem, zamiast znikac po cichu

**MessageType.CommandExtensionRegistryList** — Wyliczenie ExtensionOrigin ma dzis dwie wartosci, wiec rejestr organizacji nie ma jak sie w katalogu pokazac

**MessageType.CommandExtensionUpdateCheck** — Extension niesie wlasna wersje, ale nie wersje dostepna, wiec wskaznika aktualizacji nie ma dzis z czego zlozyc

**MessageType.CommandExtensionVersionPin** — Wersja jest dzis polem opisowym, wiec nie ma czego przypiac

**MessageType.CommandExtensionVersionRollback** — Wykonuje sie od razu; ustawienie extension.rollback.confirm wlacza potwierdzenie, ktore nie warunkuje wykonania

**MessageType.CommandExtensionBundleInstall** — Wynik jest bilansem: pozycja odrzucona wraca z powodem

**MessageType.CommandExtensionHistoryList** — Rejestr niesie stan biezacy, nie droge, ktora do niego doprowadzila

**MessageType.CommandExtensionAdminBulk** — Wynik jest bilansem przyjetych i odrzuconych, nie pojedynczym potwierdzeniem

**MessageType.CommandExtensionToolList** — Odpowiedz na tools/list, resources/list i prompts/list nie ma dzis w kontrakcie ani komendy, ani struktury, w ktora mialaby wejsc

**MessageType.CommandExtensionToolCall** — Wywolanie jest probne: idzie poza kontekstem eksperta i niczego mu nie przypisuje

**MessageType.CommandExtensionTransportSet** — Transport da sie dzis wpisac wylacznie do nieprzezroczystej konfiguracji, wiec rdzen i okno moga rozumiec to pole inaczej

**MessageType.CommandExtensionCredentialBind** — Moduł operuje na referencjach; tresc poswiadczenia nigdy nie opuszcza rdzenia i nie wchodzi w to zadanie

**MessageType.CommandExtensionSecretList** — Oddaje wylacznie odwolania i metryke; tresc sekretu nie opuszcza rdzenia

**MessageType.CommandExtensionPermissionList** — Extension niesie nieprzezroczysta konfiguracje, ktora deklaracja uprawnien nie jest

**MessageType.CommandExtensionPermissionGrant** — Nadanie jest jedyna kontrola; brak nadania nie wstrzymuje instalacji ani wlaczenia (zasada zero blokad)

**MessageType.CommandExtensionManifestScan** — Wynik jest SYGNALEM, nie brama: nie wstrzymuje instalacji ani wlaczenia

**MessageType.CommandAppsStageList** — Odpowiednik odczytu dla zdarzenia apps.build.changed — bez niego wykaz etapow zaczyna pusty przy kazdym wejsciu w modul, choc rdzen trzyma etapy w bazie

**MessageType.CommandAppsStageSave** — AppStage.ownerAgentId jedzie dzis wylacznie w kierunku od rdzenia, wiec przypisania wykonawcy nie ma jak zlecic

**MessageType.CommandAppsProductLinkList** — Powiazania nie sa czynne domyslnie — kazde jest swiadoma decyzja Operatora

**MessageType.CommandAppsArchitectureValidate** — Dzis pole validationIssues jest zapisywane i oddawane bez wyliczenia, wiec jego pustka nie orzeka o poprawnosci ukladu. Zastrzezenia sa OSTRZEZENIEM, nie brama — nie wstrzymuja pracy w warsztatach

**MessageType.CommandAppsArchitectureVersionList** — AppArchitecture.version jest liczba, ale historii kontrakt nie niesie, wiec porownania nie ma z czym zestawic

**MessageType.CommandAppsPreviewStart** — Podglad na zywo jest narzedziem warstwy 1 Frontend Workspace, a rdzen nie stawia dzis serwera produktu

**MessageType.CommandAppsEndpointList** — AppComponent niesie apiContract jako JEDEN tekst, wiec punkt koncowy nie jest dzis bytem, ktory dalby sie wyliczyc

**MessageType.CommandAppsEnvironmentList** — Wyliczenie AppDeployEnvironment ma trzy wartosci stale, a opracowanie zada dowolnej liczby srodowisk konfigurowanych przez Operatora

**MessageType.CommandAppsEnvironmentVariableList** — Wartosc sekretna oddawana jest jako referencja, nigdy jako tresc

**MessageType.CommandAppsServiceLogRead** — Podglad na zywo idzie zdarzeniem apps.service.log

**MessageType.CommandAppsArtifactList** — AppDeployment niesie wersje i odnosnik logu, ale artefakt nie jest dzis bytem umowy

**MessageType.CommandAppsDeploymentLogRead** — AppDeployment niesie logRef — sam odnosnik, bez tresci — wiec podgladu logu nie ma dzis czym wypelnic. Podglad na zywo idzie zdarzeniem apps.deployment.log

**MessageType.CommandAppsPackageValidate** — Zastrzezenia sa OSTRZEZENIEM, nie brama: nie wstrzymuja publikacji

**MessageType.CommandAppsPackageSign** — Klucz wydawcy lezy w warstwie sekretow i nie wchodzi w to zadanie ani w odpowiedz

**MessageType.CommandSpeechAudioUpload** — Jest to brak platformowy, nie brak jednego okna: speech.transcribe bierze audioRef, czyli SCIEZKE PLIKU na maszynie silnika, a nagranie z mikrofonu istnieje wylacznie jako bajty w pamieci karty. Ten sam brak wstrzymuje dyktowanie w oknie komunikacji (client/src/okno-komunikacji/dyktowanie/dostarczenie-nagrania.ts) i mikrofon Voice Console — jedna komenda otwiera obie drogi naraz. Dzwiek nie opuszcza maszyny rdzenia: bajty ida do magazynu nagran, nie do sieci

**MessageType.CommandSpeechAudioFetch** — Domyka pare do speech.audio.upload i do translate.speech.synthesize: dzis rdzen oddaje odnosniki (AssistantActivityEntry.audioRef, AssistantVoiceCommandResponse.speechRef), lecz nie ma komendy, ktora by po nie siegnela — odslucha nie ma czym wykonac

**MessageType.CommandSpeechWakeGet** — Brak modelu frazy jest ODPOWIEDZIA, nie awaria — tak samo jak brak silnika w speech.availability.get

**MessageType.CommandSpeechWakeSet** — Pola pominiete zostaja bez zmian

**MessageType.CommandSpeechListenStart** — Rdzen slucha strumienia i oglasza zdarzeniami: czesciowa transkrypcje oraz wykrycie frazy wybudzajacej. Nasluch nie jest bramka — Operator zatrzymuje go speech.listen.stop, a rdzen nie zatrzymuje go sam

**MessageType.CommandSpeechListenStop** — Zatrzymanie nasluchu, ktorego nie ma, nie jest bledem

**MessageType.CommandAssistantActivityFlag** — Dzis wyroznienie wpisu nie mialoby gdzie zamieszkac: AssistantActivityEntry nie niesie takiego pola, wiec okno pokazywaloby wyroznienie, ktorego rdzen nie pamieta

**MessageType.CommandMemoryRetentionGet** — Odczyt do pary z memory.retention.set

**MessageType.CommandMemoryRetentionSet** — Zasada obejmuje zapisy kolejne, a wpisow zastanych nie rusza wstecz — inaczej zmiana reguly kasowalaby ustalenia, na ktore Operator sie nie umawial

**MessageType.CommandMemoryContextList** — Dzis memory.toggle przestawia poziomy zasiegu i to jedyny kontekst przelaczalny; nazwy zestawu nie ma gdzie odlozyc

**MessageType.CommandMemoryContextSave** — Puste contextId zaklada nowy, podane zmienia istniejacy — tak samo jak memory.set

**MessageType.CommandMemoryContextActivate** — Aktywacja nie kasuje kontekstu poprzedniego — wraca sie do niego tym samym wywolaniem

**MessageType.CommandMemoryContextDelete** — Wpisow pamieci nie kasuje — kontekst jest zestawem wskazan, a nie ustalenie budowy tresci

**MessageType.CommandToolsScopeList** — Dzis takiego zakresu nie ma dla profilu: agent.connector.* i agent.permission.set dotycza eksperta modulu Agents i nie siegaja profilu asystenta

**MessageType.CommandToolsScopeSet** — Zakres jest nastawa zasiegu, nie bramka: stanem wyjsciowym jest pelny dostep bez limitu, zgodnie z zasada zero blokad

**MessageType.CommandClipboardList** — Schowek jest bytem klienta, ktorego rdzen dzis nie zna, wiec historia ginie wraz z karta — komenda daje jej trwalosc

**MessageType.CommandClipboardPush** — Powtorzenie tresci identycznej nie mnozy wpisow — podnosi wpis zastany na czolo wykazu

**MessageType.CommandClipboardPin** — Przypiety nie wygasa wraz z zasada retencji historii

**MessageType.CommandSnippetList** — Skrot dziala we wszystkich polach tekstowych platformy, wiec slownik nalezy do rdzenia, nie do jednego okna

**MessageType.CommandSnippetSet** — Puste snippetId zaklada nowy, podane zmienia istniejacy

**MessageType.CommandLauncherHotkeyGet** — Skrot globalny rejestruje powloka programu okiennego, wiec brak wsparcia powloki jest ODPOWIEDZIA, nie awaria — klient nie obiecuje wtedy skrotu, ktory nikogo nie obudzi

**MessageType.CommandLauncherHotkeySet** — Skrot zajety przez inny program nie jest bledem zapisu — zapis zostaje, a odpowiedz mowi, ze rejestracja sie nie udala

**MessageType.CommandContextUsageGet** — ZALEZNOSC: wymaga tokenizatora jako podsystemu rdzenia — bez niego liczba tokenow bylaby wartoscia wzieta znikad. Do czasu jego zbudowania komenda ma zwracac available falsz wraz z powodem, tak jak speech.availability.get przy braku silnika

**MessageType.CommandDesignAssetContentGet** — Pole uri zasobu jest sciezka w systemie plikow rdzenia (magazyn oddaje filepath), wiec przegladarka nie wczyta spod niego niczego — take wtedy, gdy zasob powstal bez zarzutu. Komenda dotyczy KAZDEGO zasobu magazynu, nie tylko obrazu: jeden magazyn obsluguje rodziny design.*, document.*, media.* i archive.*. Rdzen odmawia zasobu, ktorego nie zna, i zasobu, ktorego tresci nie ma pod suma kontrolna — nigdy nie oddaje bajtow zastepczych

**MessageType.CommandDesignAssetExport** — Rozni sie od image.convert celem: konwersja zaklada NOWY zasob w magazynie i tam sie konczy, eksport oddaje bajty gotowe do zapisania poza produktem

**MessageType.CommandDesignAssetExportBatch** — Pokrywa eksport zbiorczy Assets Panel oraz zestawy rozmiarow kampanii; odmowa jednego zasobu NIE wstrzymuje pozostalych — wynik niesie bilans przyjetych i odrzuconych

**MessageType.CommandDesignCollectionCreate** — Etykiety zasobu juz sa (design.asset.tag.set), ale kolekcja jest bytem osobnym: ma nazwe, opis i porzadek, a etykieta jest tylko slowem

**MessageType.CommandDesignCollectionAssign** — Zestaw jest DOKLADKA albo ODJECIEM, nie zastapieniem — inaczej niz przy etykietach, bo kolekcja bywa duza i przepisywanie jej w calosci przy kazdej zmianie jest droga do zgubienia zawartosci

**MessageType.CommandDesignCollectionList** — Bez tego kolekcja zalozona nie mialaby jak wrocic na ekran po odswiezeniu

**MessageType.CommandDesignPromptTemplateSave** — Dzis historia promptow i szablony zyja w oknie do zamkniecia karty przegladarki i tyle o nich wiadomo

**MessageType.CommandDesignPromptHistoryList** — Domyka tez prowenancje: dzis zasob niesie pole promptId, ktorego rdzen NIE wypelnia, bo nie ma przekladu klucza wiersza promptu na kod kontraktu — ta komenda ten przeklad wnosi

**MessageType.CommandDesignBoardVersionRestore** — Przywrocenie ZAKLADA nowa wersje z ukladu sprzed przywrocenia, zeby cofniecie sie samo nie kasowalo stanu, ktory Operator wlasnie porzucil

**MessageType.CommandDesignBoardExport** — Dzis kompozycja jezdzi do rdzenia i z powrotem jako uklad warstw i nie ma drogi wyjscia poza rdzen

**MessageType.CommandDesignAnnotationSet** — Pole note warstwy niesie JEDNO zdanie bez autora i bez watku; opracowanie wymaga watkow i oznaczen osob, a tego jedno pole nie unosi

**MessageType.CommandDesignPresenceReport** — Rdzen rozglasza je pozostalym zdarzeniem design.board.presence. Zgloszenie jest ULOTNE — nie zapisuje sie w bazie, bo polozenie kursora sprzed godziny nie jest wiedza o niczym

**MessageType.CommandDesignTokensetSave** — Dzis Tokens & System Panel czyta zetony z motywu obowiazujacego i nie ma ich gdzie odlozyc — kontrakt nie zna bytu zestawu zetonow

**MessageType.CommandDesignTokensetExport** — Klient sklada dzis zmienne CSS, SCSS, konfiguracje Tailwind i modul JavaScript sam i nie potrzebuje do tego rdzenia; ta komenda jest potrzebna dla postaci, ktorych przegladarka zlozyc nie moze, oraz dla wydania idacego DO INNEGO MODULU zamiast do pliku

**MessageType.CommandDesignTokensetImport** — Rdzen NIE nadpisuje motywu produktu — motyw jest wlasnoscia powloki; import zaklada byt obok niego i oddaje roznice wobec zetonow wskazanego motywu

**MessageType.CommandDesignStyleguidePublish** — Przegladarka sklada przewodnik sama i oddaje go plikiem, ale wydania go do Library albo Studio nie ma czym zlecic

**MessageType.CommandImageVectorize** — Rodzina image.* zna przeksztalcenie geometryczne, poprawke, konwersje formatu, powiekszenie i wyciecie tla — zamiany rastra na sciezki nie zna zadna z nich, a jest to funkcja opracowania modulu

**MessageType.CommandImageLayersSplit** — Zasila rozdzielenie generacji na warstwy oraz zaznaczanie obiektu; bez silnika segmentacji ODMAWIA, nazywajac brak — nigdy nie oddaje calego obrazu jako jednej warstwy udajac rozklad

**MessageType.CommandImageCompose** — Pokrywa znak wodny, branding wsadowy, osadzenie w ramce urzadzenia i warstwy rastrowe — cztery funkcje opracowania, ktore wszystkie potrzebuja jednego: zlozenia obrazu na obrazie

**MessageType.CommandRoundtableDebateGet** — Obszar nie ma dzis zadnej komendy odczytu, wiec okno otwarte w trakcie debaty zna wylacznie to, co uslyszalo zdarzeniem od swojego otwarcia

**MessageType.CommandStudioDocumentFormatSet** — Dzis format nadaje rdzen przy wczytaniu i zadna komenda go nie zmienia: studio.document.save przyjmuje documentId, content, title i createVersion, ale nie format

**MessageType.CommandStudioCommentResolve** — Watek nie znika: rozwiazanie jest stanem, nie usunieciem

**MessageType.CommandStudioTrackingSet** — Przy wlaczonym sledzeniu kazdy zapis odklada wstawienia i usuniecia jako zmiany do decyzji, zamiast nadpisywac tresc

**MessageType.CommandStudioTrackingDecide** — Decyzja zapisuje sie w tresci dokumentu i zaklada wersje

**MessageType.CommandStudioAssetEmbed** — Realizuje powiazanie Design do Studio

**MessageType.CommandStudioOperationDelete** — Operacji fabrycznej nie usuwa — na nia odpowiada odmowa nazywajaca ten fakt

**MessageType.CommandStudioChainRun** — Przebieg prowadzi petla wykonawcza okna, a kazdy krok odklada wlasna propozycje zmiany

**MessageType.CommandStudioBatchRun** — Kazdy dokument dostaje wlasne zadanie petli; odmowa jednego nie wstrzymuje pozostalych

**MessageType.CommandStudioAnnotationAdd** — Dzis czynnosc ta idzie droga generyczna window.action i wraca odmowa not_found, bo katalog akcji nie ma jej wiersza

**MessageType.CommandStudioDiffReportExport** — Wynik jest zasobem magazynu rdzenia

**MessageType.CommandStudioDiffVisual** — Sluzy tam, gdzie roznica tekstowa nie widzi zmiany ukladu

**MessageType.CommandStudioSearchSemantic** — Uzupelnia wyszukiwanie wzorca w studio.diff.compare

**MessageType.CommandStudioDiffSource** — Sluzy kontroli, czy redakcja nie odeszla od zrodla

**MessageType.CommandStudioProposalDecide** — Dzis decyzja zapada wylacznie w kliencie i rdzen o niej nie wie, wiec propozycja zostaje w nim nierozstrzygnieta

**MessageType.CommandStudioBranchCreate** — Sluzy prowadzeniu dwoch redakcji obok siebie

**MessageType.CommandStudioBranchMerge** — Konflikt nierozstrzygniety wraca w wyniku zamiast byc rozstrzygniety domyslem

**MessageType.CommandStudioPreviewRender** — Zwraca strony jako zasoby, wiec podglad pokazuje UKLAD, a nie sam tekst

**MessageType.CommandStudioPdfSplit** — Kazda czesc jest osobnym zasobem

**MessageType.CommandStudioSecurityEncrypt** — Czynnosc jest jawnym, odwracalnym ustawieniem Operatora i niczego nie warunkuje

**MessageType.CommandStudioSecurityRedact** — Czynnosc jest nieodwracalna dla wyniku, dlatego zrodlo zostaje nietkniete

**MessageType.CommandStudioSecuritySensitiveDetect** — Czynnosc CZYTA i niczego nie zmienia

**MessageType.CommandStudioIngestQueueAdd** — Dzis kolejka jest wylacznie kliencka i ginie z odswiezeniem okna

**MessageType.CommandStudioIngestRecognize** — Komenda document.text.extract robi to samo waskim wejsciem — jeden jezyk, bez silnika, bez progu i bez ukladu

**MessageType.CommandStudioIngestUrl** — Migawke strony oddaje browser.snapshot.get, ale wymaga okna modulu Browser i nie prowadzi do dokumentu Studia

**MessageType.CommandStudioIngestDeviceList** — Bez tego wykazu wybor urzadzenia nie ma z czego powstac

**MessageType.CommandStudioStyleSave** — Zmiana stylu przestawia wszystkie miejsca dokumentu, ktore go uzywaja

**MessageType.CommandStudioStyleDelete** — Stylu fabrycznego nie usuwa — odpowiada odmowa nazywajaca powod

**MessageType.CommandStudioPageSetupSet** — Zmiana formatu przelicza uklad i oddaje bilans tego, co sie nie zmiescilo

**MessageType.CommandStudioTableStructureEdit** — Szerokosci kolumn zostaja policzone, nie zerowe

**MessageType.CommandStudioDocumentImportPdf** — Odzyskanie jest odtworzeniem, nie odczytem, wiec odpowiedz niesie bilans; PDF ze samych skanow kieruje na rozpoznanie tekstu

**MessageType.CommandStudioDocumentCopy** — Kopia jest osobnym dokumentem, nie drugim odwolaniem do tego samego

**MessageType.CommandStudioDocumentExportFormat** — Format ubozszy niz dokument oddaje wykaz cech pominietych

**MessageType.CommandStudioTemplateDelete** — Szablonu fabrycznego nie usuwa — odpowiada odmowa nazywajaca powod

**MessageType.CommandStudioLockAdd** — Blokada obowiazuje w rdzeniu, przed dotknieciem tresci

**MessageType.CommandStudioLockRemove** — Blokade zdejmuje WYLACZNIE Operator — czynnosc modelu wraca odmowa nazywajaca powod

**MessageType.CommandStudioJournalRevert** — Czynnosc bedaca podstawa pozniejszej odmawia i nazywa zaleznosc, zamiast zostawic dokument w stanie niespojnym

**MessageType.CommandStudioModelChangesRevert** — Nie jest to przywrocenie wersji sprzed, bo to skasowaloby prace Operatora

**MessageType.CommandStudioAutosaveRun** — Nieudany zapis wraca nazwany, nie przemilczany

**MessageType.CommandStudioViewSet** — Gdzie da sie zrobic dwojako i obie drogi maja sens, wybor nalezy do Operatora i jest jawnym, odwracalnym ustawieniem

**MessageType.CommandStudioMarkupAdd** — Znakowanie modelu jest podpisane jako model

**MessageType.CommandStudioMarkupTypeDelete** — Rodzaju fabrycznego nie usuwa — odpowiada odmowa nazywajaca powod

**MessageType.CommandStudioAgentsSettingsGet** — Oba narzedzia sa domyslnie wylaczone

**MessageType.CommandStudioAgentsSettingsSet** — Wartosci ida zasiegami rodziny config, nie osobnym magazynem

**MessageType.CommandStudioPlanCreate** — Sam rozklad niczego nie uruchamia

**MessageType.CommandStudioPlanRun** — Gdy nastawa petli jest wylaczona, wraca odmowa nazywajaca brak nastawy, a nie cisza

**MessageType.CommandStudioAgentsClaim** — Fragment zajety przez kogo innego wraca odmowa nazywajaca wykonawce i zakres

**MessageType.CommandStudioAgentsConflictsList** — Odlozone brzmienie nie przepada

**MessageType.CommandTerminalSessionClose** — Dzis zamkniecie karty zyje wylacznie w widoku klienta: powloka i jej procesy biegna dalej, a rdzen o zamknieciu nie wie

**MessageType.CommandTerminalSessionList** — Rdzen odtwarza karty przy starcie, ale klient po ponownym polaczeniu nie ma jak ich zobaczyc i zaczyna wykaz od pustego

**MessageType.CommandTerminalFileRead** — Bez tej komendy klient czyta manifest projektu poleceniem powloki, wiec wykrycie zadan zalezy od programu wypisujacego plik i od skladni kazdej z szesciu powlok

**MessageType.CommandTerminalProcessSuspend** — Dzis rdzen umie proces wylacznie zakonczyc, wiec dlugie zadanie da sie tylko ubic

**MessageType.CommandTerminalHostSave** — Bez tej komendy ksiazka zyje jedno posiedzenie przegladarki i ginie przy odswiezeniu strony

**MessageType.CommandTerminalHostRemove** — Karty juz otwarte do tego hosta biegna dalej

**MessageType.CommandTerminalScriptSave** — Bez tej komendy biblioteka zyje jedno posiedzenie, a jedyna droga jej zachowania jest wywoz do pliku

**MessageType.CommandTerminalScriptLint** — Analize prowadza programy spoza instalki Danaco Console, wiec odpowiedz mowi wprost, czy narzedzie bylo dostepne — pusty wykaz uwag przy braku narzedzia znaczylby falszywie tresc bez zastrzezen

**MessageType.CommandTerminalTunnelOpen** — Dzis tunel da sie zalozyc wylacznie poleceniem wydanym w karcie, a wtedy jego stan i przepustowosc sa niewidoczne

**MessageType.CommandTerminalKeyGenerate** — Haslo klucza wchodzi odwolaniem do sejfu, nigdy trescia — zgodnie z zasada zapisana w kontrakcie przy zmiennej srodowiska

**MessageType.CommandTerminalKeyImport** — Klucz wskazuje sie sciezka, a nie trescia: material kluczowy nie ma powodu przechodzic przez lacze

**MessageType.CommandTerminalKeyRemove** — Wpisy ksiazki hostow wskazujace ten klucz traca wskazanie i wracaja do klucza domyslnego konfiguracji maszyny

**MessageType.CommandTerminalWatchStart** — Kontrakt daje dzis wyzwalacz plikowy automatyce, a nie karcie powloki

**MessageType.CommandTerminalWatchStop** — Polecenie juz uruchomione biegnie dalej

**MessageType.CommandWorkspaceAgentUnassign** — Nie usuwa eksperta z biblioteki modulu Agents — znosi wylacznie jego przypisanie do tego projektu

**MessageType.CommandWorkspaceProjectStatusSet** — Droga do stanu 'paused', ktorego opracowanie wymaga, a ktorego zaden dzisiejszy uchwyt nie potrafi zapisac

**MessageType.CommandWorkspaceInstructionsVersionRestore** — Przywrocenie zaklada wersje nowa o tresci wersji wskazanej — historia nie jest przepisywana

**MessageType.CommandWorkspaceTaskUpdate** — Pola pominiete zostaja bez zmiany — wywolanie nie jest podmiana calego zadania

**MessageType.CommandWorkspaceTaskMove** — Osobno od workspace.task.update, bo przeciagniecie karty zmienia stan i porzadek naraz, a jest czynnoscia jednym ruchem myszy

**MessageType.CommandWorkspaceTaskDependencySet** — Rdzen odmawia zalozenia zaleznosci domykajacej cykl, bo cykl nie da sie ulozyc w czasie

**MessageType.CommandWorkspaceScheduleGet** — Zadanie bez obu granic czasu do harmonogramu nie wchodzi

**MessageType.CommandWorkspaceNoteSave** — Puste 'noteId' zaklada notatke nowa; podane zmienia istniejaca. Rdzen przy zapisie przelicza odnosniki tresci

**MessageType.CommandWorkspaceKnowledgeGraphGet** — Nazwa rodziny 'workspace.knowledge' nie miesza sie z rodzina 'knowledge.*': tamta prowadzi wskaznik ZNACZENIA, ta rysuje siec ODNOSNIKOW miedzy bytami projektu

**MessageType.CommandWorkspaceLibraryTextExtract** — Obejmuje rozpoznanie tekstu z obrazow i skanow

**MessageType.CommandWorkspaceLibraryDuplicateList** — Sam wykaz niczego nie scala

**MessageType.CommandWorkspaceLibraryDuplicateMerge** — Etykiety, kolekcje i wersje plikow scalanych przechodza na plik zachowany, zeby scalenie nie gubilo dorobku

**MessageType.CommandWorkspaceSearchProject** — Nie zastepuje library.file.search: tamta komenda przeszukuje biblioteke centralna, ta jeden projekt i wiecej niz pliki. Wyszukiwanie po ZNACZENIU prowadzi rodzina knowledge.

**MessageType.CommandLibraryMetadataGet** — Metadane osadzone czyta sie z bajtow, wiec ich odczyt jest kosztowny i wchodzi wylacznie na wyrazne zadanie

**MessageType.CommandLibraryMetadataSet** — Dzis formularz opisu nie ma dokad pojsc: library.tag.set zmienia wylacznie etykiety i kolekcje. Zapis rozglasza library.file.changed

**MessageType.CommandLibrarySchemaSet** — Pole o kodzie juz istniejacym jest zmieniane, nie dublowane

**MessageType.CommandLibraryTagList** — Bez tej komendy klient sklada slownik z etykiet plikow odczytanej strony wykazu, wiec pokazuje probke zamiast slownika

**MessageType.CommandLibraryTagUpdate** — Zmiana nazwy przechodzi na wszystkie zasoby noszace etykiete — inaczej powstalaby druga etykieta o tym samym znaczeniu

**MessageType.CommandLibraryTagMerge** — Zasoby noszace etykiete zrodlowa dostaja docelowa, a zrodlowa znika ze slownika

**MessageType.CommandLibraryTagRemove** — Usuniecie etykiety uzywanej zada potwierdzenia, bo zdejmuje ja z zasobow, ktorych zadanie nie wymienia

**MessageType.CommandLibraryCollectionList** — Bez tej komendy klient zna wylacznie identyfikatory kolekcji wyczytane z plikow, a nazwe kolekcji tylko w chwili jej zalozenia

**MessageType.CommandLibraryRuleSet** — Zapis uruchamia przeliczenie

**MessageType.CommandLibraryRuleRemove** — Zasoby przypisane przez regule kolekcji inteligentnej zostaja w kolekcji, ale przestaja byc przeliczane — inaczej usuniecie reguly oproznialoby kolekcje po cichu

**MessageType.CommandLibraryDuplicateScan** — Rozpoznanie przyblizone zada porownania calego zbioru, wiec jest czynnoscia rdzenia, nie okna

**MessageType.CommandLibraryDuplicateMerge** — Zasoby wchloniete trafiaja do archiwum, nie znikaja

**MessageType.CommandLibraryFixityCheck** — Klient nie ma czym tego zrobic, bo nie siega po bajty zasobu

**MessageType.CommandLibraryNameNormalize** — Przebieg probny pokazuje wynik bez zapisu

**MessageType.CommandLibraryAuditList** — Dziennik jest przyrostowy — wpisow nie da sie zmienic ani usunac ta rodzina komend

**MessageType.CommandLibraryRetentionSet** — Twarde usuniecie nie zachodzi samoczynnie: polityka najwyzej zglasza zasob do usuniecia

**MessageType.CommandLibraryFileMove** — Przeniesienie nie rusza tresci ani historii wersji

**MessageType.CommandLibraryFileArchive** — Czynnosc jest odwracalna komenda library.file.restore — to jest kosz repozytorium, nie usuniecie

**MessageType.CommandLibraryFileDelete** — Jedyna droga utraty danych repozytorium; bez potwierdzenia jest odmowa

**MessageType.CommandLibraryShareCreate** — Token jest jawny zgodnie z zasada jawnosci kluczy platformy

**MessageType.CommandLibraryShareRevoke** — Odnosnik przestaje dzialac, wpis zostaje w dzienniku audytu

**MessageType.CommandLibraryClassifyRun** — Wynik wraca sugestiami do przyjecia, nie zapisem — domyslnym zachowaniem modulu jest sugestia z akceptacja Operatora

**MessageType.CommandLibrarySuggestionApply** — Przyjecie wykonuje czynnosc, ktora sugestia opisuje; odrzucenie zdejmuje ja z wykazu

**MessageType.CommandLibraryDiffCompare** — Dokumenty binarne porownuje po tekscie z nich wydobytym i mowi o tym w odpowiedzi

**MessageType.CommandDesignVectorPathSet** — Wezly przychodza w calosci — zmiana jednego wezla idzie ta sama droga co narysowanie sciezki, zeby nie bylo dwoch prawd o jej ksztalcie

**MessageType.CommandDesignVectorShapeAdd** — Ksztalt powstaje od razu jako wezly, wiec da sie go dalej edytowac pioram — nie jest osobnym bytem, ktory potem trzeba zamieniac

**MessageType.CommandDesignVectorBoolean** — Sciezki zrodlowe znikaja albo zostaja wedle wskazania — bo suma dwoch ksztaltow bywa krokiem posrednim, a bywa wynikiem koncowym

**MessageType.CommandDesignVectorTextPath** — Zamiana w kontury jest nieodwracalna dla wyniku, dlatego tekst zrodlowy zostaje

**MessageType.CommandDesignVectorOptimize** — Odpowiedz niesie ubytek zmierzony, zeby Operator widzial, ile naprawde ubylo, zamiast czytac obietnice

**MessageType.CommandDesignVectorSymbolSet** — Zmiana definicji propaguje do wszystkich instancji — po to symbol jest

**MessageType.CommandDesignFrameSet** — Ramka jest ekranem makiety: to ona, a nie kanwa, wyznacza obszar wydania

**MessageType.CommandDesignFrameRemove** — Warstwy ramki zostaja na kanwie — usuniecie ramki nie jest usunieciem pracy, ktora w niej lezala

**MessageType.CommandDesignLayoutAuto** — Rdzen oddaje warstwy po przeliczeniu, wiec klient nie liczy ukladu drugi raz i nie ma jak sie rozjechac

**MessageType.CommandDesignFrameResizeApply** — Osobno od design.frame.set, bo tam zmiana rozmiaru jest zapisem nastawy, a tu przeliczeniem ukladu wedle wiezi

**MessageType.CommandDesignMockupGenerate** — Wynik jest makieta do poprawienia, nie projektem koncowym, i tak sie o nim mowi

**MessageType.CommandDesignMockupImport** — Odpowiedz mowi, czego rdzen nie rozpoznal — bilans zamiast ciszy

**MessageType.CommandDesignColorVisionSimulate** — Wynik jest nowym zasobem — zrodlo zostaje nietkniete

**MessageType.CommandDesignColorAccessibilityAudit** — Czynnosc CZYTA i niczego nie zmienia

**MessageType.CommandDesignStockSearch** — Brak skonfigurowanego dostawcy jest ODPOWIEDZIA nazywajaca brak, nie cisza

**MessageType.CommandDesignPrintPreflight** — Czynnosc CZYTA i niczego nie zmienia; zastrzezenie z waga blad zatrzyma druk, wiec Operator ma je zobaczyc TUTAJ, a nie w drukarni

**MessageType.CommandDesignLargeformatTile** — Odpowiedz niesie uklad kafli, wiec Operator wie, ktory kafel gdzie idzie, zanim cokolwiek wydrukuje

**MessageType.CommandDesignChartRender** — Dane przychodza seriami, a nie obrazem — wykres da sie wiec przerysowac po zmianie liczb, zamiast rysowac go od nowa

**MessageType.CommandDesignPhotoCrop** — Wynik jest wariantem zrodla — oryginal zostaje nietkniety

**MessageType.CommandDesignPhotoUpscale** — Gdy kanal modelu obrazowego stoi, liczy kanal; gdy nie stoi, liczy rachunek wkompilowany — odpowiedz mowi to polem computedBy

**MessageType.CommandDesignPhotoInpaint** — Gdy kanal modelu obrazowego stoi, liczy kanal; gdy nie stoi, liczy rachunek wkompilowany — odpowiedz mowi to polem computedBy

**MessageType.CommandDesignPhotoExpand** — Gdy kanal modelu obrazowego stoi, liczy kanal; gdy nie stoi, liczy rachunek wkompilowany — odpowiedz mowi to polem computedBy

**MessageType.CommandDesignPhotoBackgroundRemove** — Gdy kanal modelu obrazowego stoi, liczy kanal; gdy nie stoi, liczy rachunek wkompilowany — odpowiedz mowi to polem computedBy

**MessageType.CommandAodMuteGet** — Wyciszenie przeterminowane nie wchodzi do wykazu — Operator nie ma go odklikiwac

**MessageType.CommandAodMuteSet** — Zniesienie idzie ta sama komenda z polem muted rownym falszowi — odwracalnosc jednym ruchem jest wymogiem, nie wygoda. Odpowiedz oddaje wykaz wyciszen PO zmianie, zeby powloka nie musiala pytac drugi raz. Zmiana rozglasza sie zdarzeniem aod.mute.changed na pozostale powloki Operatora

**MessageType.CommandAodSignalReport** — Nosnik dla trzech klas rozdz. 3.2, ktorych rdzen nie widzi wlasna telemetria: wyniku kontroli jakosci, harmonogramu przebiegow automatyk i powtarzalnosci czynnosci Operatora. Odpowiedz mowi wprost, czy sygnal wpadl w wyciszenie i ktore — sygnal wyciszony odklada sie nadal, bo wyciszenie wstrzymuje UJAWNIENIE, a nie zapis

**MessageType.CommandAodSignalList** — Odpowiedz nazywa wyciszenia, ktore sygnaly wstrzymaly, i liczbe wstrzymanych — cisza, po ktorej Operator nie wie, ze cos milczy, jest gorsza od braku wyciszenia

**MessageType.CommandMemoryDisableList** — Odczyt do pary z memory.disable.set; okno konfiguracji buduje z niego zakres ustawien „Pamiec"

**MessageType.CommandMemoryDisableSet** — Zniesienie idzie ta sama komenda z polem disabled rownym falszowi: odwracalnosc jednym ruchem. WYLACZENIE NIE KASUJE TRESCI — tym rozni sie od memory.delete; wpis wylaczony zostaje na miejscu i wraca w calosci po zniesieniu. Odpowiedz oddaje wykaz wylaczen PO zmianie

**MessageType.CommandNotificationList** — Filtry zawezaja wykaz do klas, wag, stanow albo zrodla; pominiete zwracaja rejestr niezamkniety. Licznik zdarzen nowych idzie zawsze, niezaleznie od filtru — plakietka paska kontekstu liczy caly rejestr, nie widok

**MessageType.CommandNotificationAcknowledge** — Wykaz pusty oznacza wszystkie zdarzenia nowe, czyli dzialanie zbiorcze centrum

**MessageType.CommandNotificationResolve** — Dzialanie, ktore je zamyka, wykonuje rodzina wlasciwa — ta komenda zapisuje wylacznie skutek w rejestrze centrum

**MessageType.CommandNotificationSnooze** — Po niej zdarzenie wraca do stanu nowe i znow liczy sie do plakietki — odlozenie nie jest zamknieciem

**MessageType.EventAuthChanged** — Sekcja Uwierzytelnianie Okna Ustawien odswieza sie bez odpytywania

**MessageType.EventDeviceChanged** — Rozglaszane do WSZYSTKICH polaczonych urzadzen, zeby uniewaznienie jednego bylo natychmiast widoczne na pozostalych ekranach — tej samej drogi uzywa odzyskanie konta, ktore uniewaznia tokeny wydane przed zmiana hasla

**MessageType.EventAdvisorConsulted** — Zdarzenie istnieje dla JAWNOSCI: rada, ktorej nie widac w strumieniu, jest przeslanka, ktorej w aktach nie ma — a agent na niej dziala. Niesie okno, doradce, podstawe doboru i skrot rady; pelna tresc wraca wynikiem komendy advisor.consult i zapisem dziennika konsultacji

**MessageType.EventAppsWorkspaceChanged** — Bez tego zdarzenia apps.workspace.update byl jedyna komenda zapisu modulu, ktora nic nie rozglaszala — drugie okno nie dowiadywalo sie o zmianie

**MessageType.EventSessionToolAttached** — Skoro model dostaje narzedzie, ktorego nie mial, Operator MUSI to zobaczyc — kazde posuniecie widac na ekranie

**MessageType.EventSessionToolDetached** — Widocznosc obowiazuje w OBIE strony: zestaw, ktory urosl na oczach Operatora, nie moze skurczyc sie po cichu — bez tego zdarzenia sasiednie okno pokazywaloby narzedzie, ktorego model juz nie ma

**MessageType.EventAppsPreviewChanged** — Bez tego zdarzenia podglad na zywo bylby odswiezaniem na zadanie, czyli nie podgladem na zywo

**MessageType.EventAppsDeploymentLog** — Osobne od apps.build.changed, ktore niesie STAN wdrozenia, a nie jego dziennik

**MessageType.EventSpeechListenPartial** — Tekst jest nietrwaly i bywa poprawiany kolejnym zdarzeniem — trwaly zapis daje dopiero speech.transcribe

**MessageType.EventAutomationExecutionLogged** — Nazwa rozna od komendy automation.execution.log, bo kontrakt nie dopuszcza tej samej nazwy w obu wykazach

**MessageType.EventDesignBoardChanged** — Dzis kompozycja jezdzi w obie strony komendami, ale drugie okno i drugie polaczenie nie dowiaduja sie o zmianie niczym — zdarzenia obszaru jest dokladnie jedno i dotyczy zasobu

**MessageType.EventDesignBoardPresence** — Zdarzenie ULOTNE — nie ma odpowiednika w bazie i nie jest odtwarzane po ponownym polaczeniu

**MessageType.EventAlertTriggered** — Zdarzenie jest droga alertu do Always On Display i do powiadomienia w aplikacji; rejestr odpytywany komenda alert.trigger.list sam by tam nie dotarl, bo alert ma znaczenie wtedy, gdy Operator nie patrzy

**MessageType.EventAodMuteChanged** — Rozgloszenie sciga cisza pozostale powloki Operatora: bez niego Operator wyciszalby w jednej, a sugestie wchodzilyby w drugiej

**MessageType.EventMemoryDisableChanged** — Rozgloszenie idzie na pozostale powloki Operatora, zeby wykaz pamieci w kazdej z nich pokazywal ten sam stan wlaczenia

**MessageType.EventNotificationRaised** — Rozglaszane do wszystkich polaczen Operatora, zeby plakietka i kolumna centrum nie musialy odpytywac rdzenia w tle; to samo zdarzenie jest podstawa Toastu w chwili wystapienia

**MessageType.EventNotificationChanged** — Rozglaszane do wszystkich polaczen, zeby obsluga na jednym urzadzeniu byla natychmiast widoczna na pozostalych

**MessageType.EventConnectionUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventHomeUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventEnvironmentUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventModuleUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventWorkspaceUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventSessionUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventWindowUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventMessageUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventConfigUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventSettingsUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventAccessUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventAccountUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventIdentityUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventChannelUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventQueueUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventContextUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventActionUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventStreamUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventProgressUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventStudioUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventAutomationUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventBrowserUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventResearchUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventLibraryUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventTranslateUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventRoundtableUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventDesignUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventAssistantUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventTerminalUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventDeveloperUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventDiagnosticsUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventAppsUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventAgentUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventMemoryUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventOrchestrationUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventIsolationUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventComponentUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventExtensionUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventRoleUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventSubagentUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventAdvisorUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventMonitorUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventScheduleUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventModelUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventMobileUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventAodUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventAuthUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventSpeechUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventToolsUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventArchiveUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventDeviceUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventDocumentUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventHistoryUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventImageUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventKnowledgeUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventMailUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventMediaUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventPanelUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventRetentionUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventTeamUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventAlertUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventClipboardUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventHealthUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventLauncherUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventProvenanceUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventSnippetUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventUsageUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventNotificationUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**Envelope** — Pola odpowiedzi i strumienia sa opcjonalne.

**Window.AgentId** — Agent nie zastepuje kanalu — kanal jest droga do modelu, agent tozsamoscia nalozona na te droge

**Window.CoordinatorWindowId** — Odpowiednik kolumny okno_komunikacji.okno_koordynatora_id

**Channel.AccountId** — Kanal mowi CZYM sie polaczyc, konto — CZYIM poswiadczeniem

**Environment** — Odpowiednik wiersza tabeli srodowisko

**Environment.Motto** — Brak motta jest stanem normalnym: puste pole ma zostac puste, nie wolno go zastapic tekstem ulozonym po stronie widoku

**Environment.ModuleCodes** — Puste dla srodowiska z panelem orkiestracji

**Module** — Odpowiednik wiersza tabeli modul

**Module.OperationalWindowCodes** — Okno rozmowy wystepuje w kazdym module

**Module.Kind** — Odpowiednik kolumny modul.rodzaj (migracja 076)

**Module.ConfiguredOnHome** — NIEZALEZNE OD `kind`: rodzaj mowi CZYM modul jest, to pole GDZIE sie go sklada

**ToolDeclaration** — Brak osobnego parsera intencji: polecenie wpisane w oknie rozmowy model wykonuje wywolaniem narzedzia

**SettingDefinition** — Niesie komplet metadanych potrzebnych do zbudowania pola formularza — klient nie zaszywa u siebie ani jednego klucza

**AccessPoint** — Punkt mowi, DO CZEGO model ma wglad; katalog roboczy modelu jest ustawieniem osobnym (klucz katalog.roboczy.podstawa)

**AccessGrant** — Okno ma ZBIOR nadan; kolejnosc i oznaczenie glownego maja znaczenie

**Account** — Poswiadczenia wchodza zadaniem i nie wychodza nigdy — odpowiedz niesie wylacznie znacznik hasCredential

**IdentityCategory** — Katalog obejmuje konstytucje, profil roli, ekspertyze, zasady bezpieczenstwa i zasady harnessu

**LoopState** — Bieg nie ma limitu obiegow — zamiast bramy licznikowej wchodzi przejrzystosc: licznik obiegow, obiegi bez postepu i jawny powod zatrzymania. Odpowiednik session.StanObiegu

**SessionPresence** — Zasila kontrolke powrotu na stronie glownej: rozlaczenie klienta nie konczy sesji ani procesow, wiec strona glowna musi wiedziec, ze sesja zyje i gdzie zyje. session.bind jest akcja powrotu wywolywana z tej kontrolki, nie ruchem otwierajacym

**SessionConfig** — NIE jest zbiorem flag programu: adapter dostawcy dopiero tlumaczy go na wiersz argumentow, zmienne srodowiskowe, plik ustawien, plik serwerow MCP, plik pamieci, pliki promptu i katalog roboczy. Ten sam model obsluguje kanal CLI, API, SDK, SSH i model lokalny — zmiana dostawcy zmienia adapter, nie model. Osiemnascie obszarow najwyzszego poziomu; warstwy tozsamosci, profilu roli i ekspertyzy sa polami obszaru systemPrompt, nie osobnymi obszarami. Kazdy obszar jest opcjonalny: brak obszaru znaczy wartosc odziedziczona z poziomu szerszego albo domyslna

**SessionConfigAccount** — Zrodlem prawdy dla konta i poswiadczenia jest rejestr kont (account.*) — konfiguracja sesji odwoluje sie do niego identyfikatorem i nigdy nie niesie tresci poswiadczenia

**SessionConfigProvider** — Zrodlem prawdy dla parametrow polaczenia jest rejestr kanalow modelu (channel.*); pola ponizej sa nadpisaniem per sesja, nie druga kopia rejestru

**SessionConfigSystemPrompt** — Tozsamosc jest ZAMIENIANA, nie dolaczana — domyslnym trybem jest ZASTAP. Warstwy odpowiadaja injection.WarstwaKonstytucja, injection.WarstwaProfil i injection.WarstwaEkspertyza; obszar STERUJE istniejacym silnikiem nakladki, nie powtarza go

**SessionConfigSystemPrompt.Mode** — Puste znaczy ZASTAP

**SystemPromptLayer** — Edycja tresci idzie WYLACZNIE komenda identity.document.set — konfiguracja sesji wskazuje kategorie, nie zastepuje katalogu tozsamosci

**SessionConfigTools** — Wykaz narzedzi platformy powstaje z kontraktu — tutaj rozstrzyga sie wylacznie ich dopuszczenie

**SessionConfigPermissions** — Pola bypassAllPrompts i promptHandlerTool sa punktami zastrzezonymi Operatorowi — model nie rozszerza wlasnych uprawnien

**SessionConfigMcp** — Zrodlem prawdy dla adresu i poswiadczenia mostu jest rejestr punktow dostepu (access.*) — wiazanie odwoluje sie do punktu identyfikatorem

**SessionConfigMemory** — Adapter zapisuje tresc do pliku pamieci katalogu roboczego — pole niesie tresc, nie sciezke pliku dostawcy

**SessionConfigConversationContext** — Przenoszeniem kompletu miedzy modulami steruje context.transfer

**SessionConfigWorkingDirectory** — Katalog roboczy jest ustawieniem OSOBNYM od dostepu i od srodowiska platformy. Rozejscie miedzy katalogiem ustawionym a faktycznie uzywanym pokazuje WorkingDirectoryResolution w konfiguracji obowiazujacej

**SessionConfigAdditionalDirectories** — Zrodlem prawdy dla zakresu wgladu jest rejestr nadan dostepu okna (access.grant.*) — dostep to nie katalog roboczy

**SessionConfigEnvironment** — To NIE jest srodowisko platformy — profil widocznosci modulow (TalkIn, WorkSpace, CodeStudio, MultitaskingAI) zyje w obszarze komend environment i pozostaje nietkniety

**EnvironmentVariable** — Wartosc tajna wchodzi odwolaniem do sejfu, nigdy trescia

**SessionConfigInputOutput** — Pola sa dziedzinowe, nie sa nazwami przelacznikow — adapter dopiero tlumaczy je na postac dostawcy

**SessionConfigRuntime** — Uruchomieniem steruje injection.Rozruch — jedyny spawner platformy

**SessionConfigSessionLifecycle** — Bieg naprawczy nie ma bramy licznikowej — zamiast niej stoi prog braku postepu i jawny powod zatrzymania; stan biegu niesie LoopState

**WorkingDirectoryResolution** — Degradacja ma byc widoczna w oknie, nie milczaca; sama degradacja nie wstrzymuje uruchomienia

**SessionConfigFieldCapability** — Pole bez odpowiednika u dostawcy jest wypisane JAWNIE wraz z powodem, zamiast byc pomijane po cichu — tak jak turnLimit i promptHandlerTool, ktorych zainstalowany program CLI nie zna. Zmiana wersji programu albo przejscie na API zmienia wylacznie adapter i te deklaracje

**AdapterCapabilities** — Klient pyta o nia przed zbudowaniem okna konfiguracji i wie, ktore pole dostawca zignoruje, zanim Operator je wypelni

**AgentVersion.Mode** — Pole NIEWYMAGANE i to niesie trzeci stan: brak wartosci oznacza prompt globalny, DOLACZ prompt dopisywany, ZASTAP odstepstwo jawne

**DiagnosticError.ErrorCode** — Odmowa komendy jest FAKTEM o niej samej, wiec kod nalezy do bledu, nie do jego kontekstu; kolumna `kod_bledu` trzyma go w bazie od migracji 043

**AgentLayer** — Warstwy eksperta DOPISUJA sie do promptu systemowego jako zakres uzytkownika i NIGDY go nie zastepuja — dlatego warstwa nie ma pola trybu. Prompt systemowy jest paczka built-in, obowiazujaca zawsze i jako pierwsza (rozstrzygniecie ustalenie budowy)

**Agent.Mode** — Pominiete znaczy DOLACZ: warstwy eksperta DOPISUJA sie do globalnego promptu systemowego z okna konfiguracji i ustawien. ZASTAP jest ODSTEPSTWEM od ustawien domyslnych, oznaczanym przez Operatora swiadomie — instrukcja eksperta STAJE SIE promptem systemowym, a globalny przestaje obowiazywac dla okien tego eksperta

**Agent.MemoryLevels** — Pole oddawane ZAWSZE, takze puste: LISTA PUSTA ZNACZY PAMIEC WYLACZONA W CALOSCI i jest to jedyny zapis wylaczenia ( MemoryLevel). Pominiecie pola zostawialoby Operatora z domyslem zamiast z odpowiedzia. Wartosc wyjsciowa nowego eksperta to cztery poziomy — stanem wyjsciowym platformy jest pelny dostep

**Agent.ModuleCodes** — Pole oddawane ZAWSZE, takze puste: LISTA PUSTA ZNACZY BRAK OGRANICZENIA, czyli dostepnosc we wszystkich modulach, i jest to stan wyjsciowy. Pominiecie pola zostawialoby Operatora z domyslem zamiast z odpowiedzia — ta sama zasada co przy `memoryLevels`

**Agent.SubagentLimit** — Wartosc wyjsciowa 15 — maksimum techniczne platformy. Zero znaczy Subagent Network wylaczony i jest to jedyny zapis wylaczenia, wiec `agent.subagent.set` z `enabled: false` zapisuje tu zero

**Extension.Origin** — Rozstrzyga WYLACZNIE stan wyjsciowy przy rejestracji; poza tym jest faktem do pokazania Operatorowi

**Subagent** — UWAGA: czym podagent JEST — oknem, procesem czy pozycja kolejki — nie zostalo rozstrzygniete przez ustalenie budowy

**AodStatus.ModuleId** — Do tej pozycji nakladka dochodzila modulu okreznie — przez window.list — i sugestia, ktorej okna rdzen nie znal, w wyciszenie modulu nie wpadala

**AodSuggestion.ModuleId** — Bez niego wyciszenie biezacego modulu nie ma po czym rozpoznac swojej sugestii

**AodSuggestion.EventClass** — Podstawa wyciszenia klasy zdarzen; puste znaczy klase nierozpoznana, a sugestia bez klasy w wyciszenie klasy nie wpada

**AuthMethod** — Sekret nie trafia do bazy — zostaje w niej wylacznie odwolanie; odpowiedz nie niesie ani skrotu hasla, ani skrotu PIN-u, ani materialu klucza prywatnego. PIN i Windows Hello sa WLASCIWE URZADZENIU: klucz nie opuszcza jego TPM, wiec na kazdej maszynie zaklada sie je osobno

**AuthSession** — Transportem jest gniazdo WebSocket, nie seria zadan HTTP, wiec ciasteczka nie ma — token wraca z auth.login i idzie w connection.hello. Wygasanie przesuwne, przedluzane aktywnoscia albo jawnie komenda auth.token.refresh

**Device** — Konto jest jedno, urzadzen dowolnie wiele; kazde niesie wlasny token dostepu, ktory Operator moze uniewaznic osobno

**MailAccount** — Poswiadczenia leza w sejfie rdzenia, nigdy w bazie ani w kontrakcie

**MailMessage** — Tresc jest wypelniona wylacznie przy mail.message.get — wykaz niesie same naglowki, bo ciagniecie tresci calej skrzynki byloby pobraniem archiwum przy kazdym pytaniu

**KnowledgeImageHit** — Osobna struktura od KnowledgeHit, bo obraz nie ma fragmentu tekstu do zacytowania — wskazaniem jest sam plik

**SessionTool** — Zyje W STANIE SESJI, nie w definicji eksperta: definicja pozostaje nietknieta, dolozenie przezywa rozlaczenie klienta i konczy sie wraz z sesja albo z usunieciem rozmowy

**ToolCatalogEntry** — Trzy pola nazewnicze nie sa ozdoba: nazwa pelna niesie przedrostek zrodla, nazwa skrocona jest tym, co Operator wpisuje, a opis ratuje przed cicha degradacja — nazwa :design-audit sama nie mowi nic (rozdz. 7)

**AgentVersionSnapshot** — Byt odrebny od `Agent`, bo wersja nie niesie ani identyfikatora eksperta, ani jego przypisan, ani stanu czynnosci — utrwalone jest to, co `agent.update` przyjmuje, i tylko to

**ChannelCheckResult** — Struktura wydzielona, bo ten sam ksztalt oddaje sprawdzenie wykonane recznie z okna i sprawdzenie wykonane przy podlaczeniu rozszerzenia rozdz. 13, punkt sterowania „Test polaczenia konektora")

**ChannelCredentialStatus** — Struktura celowo nie ma pola na wartosc: brak pola jest mocniejszym zabezpieczeniem niz zasada, ze pola sie nie wypelnia

**AppDependency** — AppComponent.dependsOn niesie same identyfikatory, wiec strzalki przeplywu danych nie ma dzis z czego narysowac

**ModelCallTrace** — Wiersz tabeli prowenancja_wywolanie

**ModelCallTrace.ProcessId** — To jest pole korelacyjne wspolne z LogEntry.processId i MonitorStatus.processId — po nim wpis dziennika laczy sie ze sladem

**ModelCallSpan** — Wiersz tabeli prowenancja_odcinek

**UsageAggregate** — Agregat powstaje z odczytu tabeli prowenancja_wywolanie, nie z osobnej tabeli zuzycia — dwie tabele o tych samych liczbach rozjechalyby sie przy pierwszej korekcie cennika

**UsageAggregate.CostWithoutPrice** — Wieksze od zera znaczy, ze koszt jest niepelny i okno ma to powiedziec

**AlertRule** — Wiersz tabeli alert_regula wraz z kanalami z tabeli alert_regula_kanal

**AlertTrigger** — Wiersz tabeli alert_wyzwolenie

**AlertTrigger.DeliveredChannels** — Roznica wobec regul.channels znaczy nieudane dostarczenie, nie brak wyzwolenia

**HealthProbe** — Wiersz tabeli kondycja_sonda

**HealthProbeResult** — Wiersz tabeli kondycja_wynik. Bez tej serii nie ma z czego policzyc ani procentu dostepnosci, ani budzetu bledow — dlatego wyniki sa wierszami, a nie ostatnim stanem nadpisywanym w miejscu

**HealthUptime.Samples** — Zero znaczy brak pomiaru, a nie dostepnosc zerowa — okno ma odroznic jedno od drugiego

**LibraryMetadata** — Piec pol profilu podstawowego (title, creator, subject, date, rights) wypelnia okno domyslnie; pozostale dziesiec nalezy do standardu Dublin Core Elements 1.1 i jest oplacalne raz, a nie przy kazdym rozszerzeniu profilu

**LibraryMetadata.Identifier** — DOI albo ISBN

**LibraryTechnicalMetadata** — Pola geolokalizacji sa jedynym zrodlem widoku mapy w Library Explorer — bez nich widok nie ma czego naniesc

**LibraryTag** — Dzis etykieta jest w kontrakcie samym napisem przy pliku i nie ma tozsamosci — bez niej nie da sie ani nadac jej barwy, ani polaczyc dwoch etykiet w jedna

**LibraryDuplicateGroup** — Podstawa rozpoznania wchodzi do wyniku, bo duplikat pewny i duplikat przyblizony zadaja innej decyzji Operatora

**LibraryFixityResult** — Suma oczekiwana i wyliczona stoja osobno, bo rownosc jest wnioskiem, a nie danymi

**LibraryStats** — Klient liczy dzis te wartosci z odczytanej strony wykazu, wiec orzeka o probce, a nie o repozytorium

**TerminalSshKey** — Czesc tajna nie opuszcza maszyny rdzenia i nie ma tu pola

**WorkspaceTask.Rank** — Napis, nie liczba: wstawienie karty miedzy dwie sasiednie nie ma wtedy przepisywac calej kolumny

**WorkspaceBacklink.TargetNoteId** — Puste znaczy odnosnik do strony jeszcze niezalozonej — to stan poprawny, nie usterka

**WorkspaceKnowledgeGraph.Truncated** — Bez tego pola przyciety graf wygladalby na kompletny obraz projektu

**WorkspaceCanvas.Scene** — Rdzen sceny nie rozbiera na wiersze: jej ksztalt nalezy do widoku, ktory ja rysuje

**StudioParagraphFormat.ListRestart** — Bez tego pola wznowienie rozdzielalo liste na dwie definicje, a zmiana jednej definicji przestawala przestawiac oba jej odcinki

**StudioFragmentLock.TemplateId** — Blokada wychodzaca z fromTemplate bez tego pola nie mowi, KTOREGO szablonu jest wzorcem, wiec Operator nie ma jak dojsc do jej zrodla

**StudioActor** — Pojecia agenta nie zaklada od nowa: agentId odpowiada Agent.id, subagentId odpowiada Subagent.id

**StudioAgentSettings** — Oba narzedzia sa DOMYSLNIE WYLACZONE i wlaczane jawnym, odwracalnym ustawieniem Operatora; wartosci ida zasiegami rodziny config, nie osobnym magazynem

**StudioAgentConflict** — Rdzen nie rozstrzyga, kto ma racje; niesie prawde o tym, co sie stalo

**AodMute** — Wyciszenie jest bytem rdzenia, nie stanem jednego okna: Operator wyciszajacy nakladke w jednej powloce ma miec cisze w kazdej nastepnej. Zniesienie idzie ta sama komenda aod.mute.set z polem muted rownym falszowi — jednym ruchem, bez pytania o potwierdzenie

**AodMute.ScopeName** — Cisza, przy ktorej Operator nie wie, CO milczy, jest gorsza od braku wyciszenia

**AodMute.EndsAt** — Puste znaczy wyciszenie trwajace do zniesienia reka Operatora — tak stoi wyciszenie kontekstowe i wyciszenie klasy zdarzen

**AodMute.DeviceId** — Wyciszenie obowiazuje wszystkie powloki Operatora niezaleznie od tego pola; pole mowi, skad przyszlo

**AodSignal** — Nosnik dla klas, ktorych rdzen nie widzi wlasna telemetria: wyniku kontroli jakosci, harmonogramu przebiegow automatyk i powtarzalnosci czynnosci Operatora (rozdz. 3.2 i 3.3 opracowania). Bez niego trzy z szesciu klas nie maja czego wyciszac ani z czego zrobic sugestii

**AodSignal.OccurrenceCount** — Prog powtarzalnosci czynnosci recznej z rozdz. 3.4 liczy sie z tego pola; puste znaczy jedno wystapienie

**MemoryDisable** — Wylacza sie albo pojedynczy wpis (entryId), albo caly poziom pamieci (level) — w zasiegu: calkiem, w srodowisku, w projekcie, w module, w parze modulow albo w karcie sesji. WYLACZENIE NIE KASUJE TRESCI i to jest jego roznica wobec memory.delete: wpis wylaczony nie wchodzi do kontekstu, ale zostaje na miejscu i wraca w calosci zniesieniem wylaczenia. Sterowanie stoi po stronie Operatora w konfiguracji — zakres ustawien „Pamiec"

**MemoryDisable.Level** — Puste znaczy wylaczenie pojedynczego wpisu wskazanego polem entryId. Zadanie bez obu pol nie ma czego wylaczyc i konczy sie odmowa

**MemoryDisabledEntry** — Wpis znika z kontekstu, nie z pamieci — tresc stoi nietknieta i wraca zniesieniem wylaczenia

**NotificationAction** — Wykonania NIE prowadzi rodzina notification.*: pozycja wskazuje komende rodziny wlasciwej i jej ladunek

**ConnectionHelloRequest.Token** — Puste znaczy polaczenie niezwiazane — bramka wtedy nie wie, kto po drugiej stronie

**ConnectionHelloResponse.LoginRequired** — Puste znaczy, ze rdzen nastawy nie zna. To nie jest bramka: pole niczego nie odmawia, mowi klientowi, czy okno logowania ma sie pokazac

**HomeEnterResponse.Presence** — Puste, gdy zadna sesja nie trwa

**SpeechAvailabilityGetRequest** — Brak silnika to ODPOWIEDZ, nie awaria — klient nie pokazuje mikrofonu, ktory nic nie nagra

**SpeechAvailabilityGetResponse** — Brak silnika to ODPOWIEDZ, nie awaria — klient nie pokazuje mikrofonu, ktory nic nie nagra

**SpeechAvailabilityGetResponse.SynthesisAvailable** — Osobne pole od available, bo odsluch jedzie innym lancuchem niz dyktowanie: piper z glosem .onnx albo espeak-ng, nie python z faster-whisper

**SpeechAvailabilityGetResponse.SynthesisReason** — Puste gdy odsluch gotowy

**SpeechTranscribeRequest** — Dzwiek NIE OPUSZCZA maszyny, na ktorej stoi silnik

**SpeechTranscribeResponse** — Dzwiek NIE OPUSZCZA maszyny, na ktorej stoi silnik

**SpeechTranscribeResponse.Transcript** — Pusty JEST dozwolony, gdy processed=true — cisza i szum to prawidlowy wynik pomiaru

**WorkspaceEnterRequest.WindowId** — Okno rozmowy nie znika przy zmianie modulu, tylko rekonfiguruje kontekst

**SessionFocusRequest** — Asystent przestawia ognisko OPERATOROWI polem targetClientId — inaczej jego posuniecia nie byly widoczne na ekranie

**SessionFocusRequest.TargetClientId** — Uzywa tego ASYSTENT dzialajacy za Operatora: bez tego pola przestawialby ognisko sobie, a ekran Operatora stalby w miejscu — czyli praca dzialaby sie "gdzies w tle"

**SessionFocusResponse** — Asystent przestawia ognisko OPERATOROWI polem targetClientId — inaczej jego posuniecia nie byly widoczne na ekranie

**WorkspaceAgentAssignRequest.Scope** — Zasieg szerszy czyni przypisanie obowiazujacym takze w innych projektach

**ResearchReportBuildResponse.FromModel** — Dzis klient nie ma jak odroznic streszczenia od komunikatu bledu, bo rdzen zbiera z kanalu same fragmenty tekstu i zapisuje je jako sekcje

**TranslateTargetAddRequest.ChannelId** — Operator nie mial dotad jak wskazac, ktory model tlumaczy

**DesignAssetGenerateRequest** — Tresc trafia do magazynu rdzenia pod suma kontrolna, wiec zasob nie zalezy od zadnego pliku zewnetrznego. Brak kanalu obrazowego albo brak poswiadczenia to ODMOWA NAZYWAJACA BRAK — rdzen nigdy nie zaklada zasobu bez bajtow obrazu

**DesignAssetGenerateRequest.ChannelId** — Brak bierze pierwszy czynny kanal obrazowy konta; brak takiego kanalu to odmowa nazywajaca brak, nigdy obraz zastepczy

**DesignAssetGenerateResponse** — Tresc trafia do magazynu rdzenia pod suma kontrolna, wiec zasob nie zalezy od zadnego pliku zewnetrznego. Brak kanalu obrazowego albo brak poswiadczenia to ODMOWA NAZYWAJACA BRAK — rdzen nigdy nie zaklada zasobu bez bajtow obrazu

**AgentCreateRequest.MemoryLevels** — Pominiete znaczy "bez zmiany" (przy zakladaniu: cztery poziomy). LISTA PUSTA `[]` JEST ZADANIEM WYLACZENIA PAMIECI i rozni sie od pominiecia pola — dlatego wylaczenie nie potrzebuje wlasnej wartosci wyliczenia

**AgentUpdateRequest.MemoryLevels** — Pominiete znaczy "bez zmiany" (przy zakladaniu: cztery poziomy). LISTA PUSTA `[]` JEST ZADANIEM WYLACZENIA PAMIECI i rozni sie od pominiecia pola — dlatego wylaczenie nie potrzebuje wlasnej wartosci wyliczenia

**AgentListRequest.ProjectId** — Podany zawezaja wykaz do ekspertow widocznych w tym projekcie: wszystkich `global` oraz tych `project`, ktore sa do niego przypisane. Pominiety znaczy wykaz biblioteki w calosci — okno Agent Buildera musi widziec takze ekspertow projektowych, inaczej nie daloby sie ich poprawic

**AgentLayerSetRequest** — Warstwa DOPISUJE sie do promptu systemowego, nie zastepuje go — trybu podania sie tu nie wybiera

**AgentLayerSetResponse** — Warstwa DOPISUJE sie do promptu systemowego, nie zastepuje go — trybu podania sie tu nie wybiera

**TranslateBacktranslationRunRequest.ChannelId** — Operator nie mial dotad jak wskazac, ktory model tlumaczy

**MemoryListResponse.DisabledEntries** — Do kontekstu nie wchodza i nie stoja w polu entries; wykaz jest zawsze tablica, bo brak wylaczen to fakt, nie cisza

**MemoryDetachRequest** — Odpiecie NIE JEST wylaczeniem: zweza zasieg samego wpisu, a nie wstrzymuje go w cudzym projekcie ani module — do tego sluzy memory.disable.set

**MemoryDetachResponse** — Odpiecie NIE JEST wylaczeniem: zweza zasieg samego wpisu, a nie wstrzymuje go w cudzym projekcie ani module — do tego sluzy memory.disable.set

**ComponentAssignRequest** — UWAGA: znaczenie przypisania nie zostalo rozstrzygniete przez ustalenie budowy

**ComponentAssignResponse** — UWAGA: znaczenie przypisania nie zostalo rozstrzygniete przez ustalenie budowy

**ExtensionInstallRequest** — UWAGA: znaczenie instalacji — pobranie paczki, zarejestrowanie adresu czy zapis punktu dostepu — nie zostalo rozstrzygniete przez ustalenie budowy

**ExtensionInstallRequest.Origin** — Rozstrzyga stan wyjsciowy: `danaco` staje wlaczone, `personal` wylaczone (rozdz. 5.2)

**ExtensionInstallResponse** — UWAGA: znaczenie instalacji — pobranie paczki, zarejestrowanie adresu czy zapis punktu dostepu — nie zostalo rozstrzygniete przez ustalenie budowy

**ExtensionUninstallRequest** — UWAGA: znaczenie odinstalowania jest zwiazane z nierozstrzygnietym znaczeniem instalacji

**ExtensionUninstallResponse** — UWAGA: znaczenie odinstalowania jest zwiazane z nierozstrzygnietym znaczeniem instalacji

**AdvisorConsultRequest** — Rada NIE JEST WIAZACA i nie wykonuje pracy za pytajacego — odpowiedzialnosc za wynik zostaje przy agencie, ktory pyta. Kim jest doradca, rozstrzyga RDZEN z danych okna, nie zadanie: pole requestedAdvisor jest PROSBA i podlega sufitowi sily. Prosba o model silniejszy od kanalu pytajacego bez wczesniejszego wskazania Operatora jest ODMAWIANA, a nie zamieniana po cichu na slabszego doradce — model z wlasnej inicjatywy po model silniejszy nie siega (rozstrzygniecie ustalenie budowy z 14.08.2026). Kazda odbyta konsultacja rozglasza sie zdarzeniem advisor.consulted

**AdvisorConsultRequest.WindowId** — Z NIEGO rdzen bierze kanal pytajacego i wskazanie Operatora — zadanie tych dwoch rzeczy nie podaje i podac nie moze

**AdvisorConsultResponse** — Rada NIE JEST WIAZACA i nie wykonuje pracy za pytajacego — odpowiedzialnosc za wynik zostaje przy agencie, ktory pyta. Kim jest doradca, rozstrzyga RDZEN z danych okna, nie zadanie: pole requestedAdvisor jest PROSBA i podlega sufitowi sily. Prosba o model silniejszy od kanalu pytajacego bez wczesniejszego wskazania Operatora jest ODMAWIANA, a nie zamieniana po cichu na slabszego doradce — model z wlasnej inicjatywy po model silniejszy nie siega (rozstrzygniecie ustalenie budowy z 14.08.2026). Kazda odbyta konsultacja rozglasza sie zdarzeniem advisor.consulted

**AdvisorConsultResponse.AdvisorChannel** — Rada nigdy nie wraca sama — bez wskazania doradcy dalaby sie podac za odpowiedz wlasna agenta

**ConfigExplainGetRequest** — Transparentnosc, nie bramka

**ConfigExplainGetResponse** — Transparentnosc, nie bramka

**ConfigWindowOpenRequest** — UWAGA: wykaz okien obiecuje te komende, nie mowiac, co rdzen ma przy niej robic poza podaniem zakresu

**ConfigWindowOpenResponse** — UWAGA: wykaz okien obiecuje te komende, nie mowiac, co rdzen ma przy niej robic poza podaniem zakresu

**TerminalOutputReadRequest** — Komenda `terminal.command.exec` konczy sie w chwili STARTU procesu (kompilacja trwa dluzej niz kazde sensowne oczekiwanie na odpowiedz, a rozlaczenie klienta nie ma prawa jej przerwac), wiec wyniku niesc nie moze i nigdy nie bedzie mogla. Tedy model czyta, co polecenie wypisalo. Zrodlem jest ten sam dziennik zbiorczego wyjscia, na ktorym stoi `terminal.output.stream` — drugiej pompy nie ma. Historia zyje jeden bieg rdzenia: po restarcie wyjscie jest puste i odpowiedz mowi to wprost

**TerminalOutputReadRequest.WaitMs** — Gorna granica to 60000 — dluzsze czekanie trzymaloby zadanie gniazda na czas, ktorego zadne polaczenie nie przezyje. Czekanie NIE jest warunkiem odpowiedzi: proces wciaz biegnacy oddaje wyjscie dotychczasowe i stan running

**TerminalOutputReadResponse** — Komenda `terminal.command.exec` konczy sie w chwili STARTU procesu (kompilacja trwa dluzej niz kazde sensowne oczekiwanie na odpowiedz, a rozlaczenie klienta nie ma prawa jej przerwac), wiec wyniku niesc nie moze i nigdy nie bedzie mogla. Tedy model czyta, co polecenie wypisalo. Zrodlem jest ten sam dziennik zbiorczego wyjscia, na ktorym stoi `terminal.output.stream` — drugiej pompy nie ma. Historia zyje jeden bieg rdzenia: po restarcie wyjscie jest puste i odpowiedz mowi to wprost

**TerminalOutputReadResponse.Truncated** — Prawda znaczy, ze poczatek zostal odciety — ostatnie bajty sa zawsze cale

**ModelChannelSetRequest** — Kanaly zaklada i zmienia rodzina channel.*; ta komenda wybiera jeden z zalozonych

**ModelChannelSetResponse** — Kanaly zaklada i zmienia rodzina channel.*; ta komenda wybiera jeden z zalozonych

**AodSuggestionRequest** — UWAGA: wykaz okien obiecuje te pozycje slowem podpowiedz, nie rozstrzygajac, czy jest to odczyt Operatora, czy zdarzenie wypychane przez rdzen

**AodSuggestionResponse** — UWAGA: wykaz okien obiecuje te pozycje slowem podpowiedz, nie rozstrzygajac, czy jest to odczyt Operatora, czy zdarzenie wypychane przez rdzen

**AuthRegisterRequest** — Konto powstaje w stanie NIEPOTWIERDZONYM i pozostaje w nim do chwili potwierdzenia adresu komenda auth.verify — dopiero potwierdzenie wydaje urzadzeniu token dostepu. Sesji ta komenda NIE zaklada. Wykonalna tylko raz; potem odmawia trwale. Powtorzenie hasla jest sprawa formularza klienta, nie kontraktu

**AuthRegisterRequest.Email** — Weryfikowany przy rejestracji, pelni pozniej funkcje metody logowania i JEDYNEJ drogi odzyskania konta — drugiego adresu do tych celow platforma nie prowadzi

**AuthRegisterResponse** — Konto powstaje w stanie NIEPOTWIERDZONYM i pozostaje w nim do chwili potwierdzenia adresu komenda auth.verify — dopiero potwierdzenie wydaje urzadzeniu token dostepu. Sesji ta komenda NIE zaklada. Wykonalna tylko raz; potem odmawia trwale. Powtorzenie hasla jest sprawa formularza klienta, nie kontraktu

**AuthRegisterResponse.PendingVerification** — Prawda oznacza, ze list z droga potwierdzenia zostal wyslany na podany adres

**AuthVerifyRequest** — Przenosi konto ze stanu niepotwierdzonego do potwierdzonego, po czym wydaje urzadzeniu token dostepu — to jest moment, w ktorym Operator wchodzi do platformy po raz pierwszy. Droga potwierdzenia przychodzi listem na adres podany przy rejestracji i wygasa; wygasla droga odmawia i pozwala poprosic o nowa

**AuthVerifyRequest.KeepSignedIn** — Nie jest bramka: znosi powtarzanie logowania, niczego nie blokuje

**AuthVerifyResponse** — Przenosi konto ze stanu niepotwierdzonego do potwierdzonego, po czym wydaje urzadzeniu token dostepu — to jest moment, w ktorym Operator wchodzi do platformy po raz pierwszy. Droga potwierdzenia przychodzi listem na adres podany przy rejestracji i wygasa; wygasla droga odmawia i pozwala poprosic o nowa

**AuthRecoverRequest** — Serwer wysyla na wskazany adres droge potwierdzenia tozsamosci; nowe haslo ustawia sie komenda auth.reset. Odpowiedz nie zdradza, czy adres pasuje do konta — inaczej komenda mowilaby obcemu, jaki adres ma Operator

**AuthRecoverResponse** — Serwer wysyla na wskazany adres droge potwierdzenia tozsamosci; nowe haslo ustawia sie komenda auth.reset. Odpowiedz nie zdradza, czy adres pasuje do konta — inaczej komenda mowilaby obcemu, jaki adres ma Operator

**AuthRecoverResponse.Sent** — Prawda niezaleznie od tego, czy adres pasuje do konta — wartosc nie jest odpowiedzia na pytanie o istnienie konta

**AuthResetRequest** — Zastepuje dotychczasowy skrot hasla przechowywany poza baza danych i uniewaznia tokeny dostepu wydane przed odzyskaniem — kazde powiazane urzadzenie loguje sie ponownie, a zmiana idzie w swiat zdarzeniem device.changed. Nowej encji Konto nie tworzy: zmienia sie wylacznie material uwierzytelniajacy

**AuthResetResponse** — Zastepuje dotychczasowy skrot hasla przechowywany poza baza danych i uniewaznia tokeny dostepu wydane przed odzyskaniem — kazde powiazane urzadzenie loguje sie ponownie, a zmiana idzie w swiat zdarzeniem device.changed. Nowej encji Konto nie tworzy: zmienia sie wylacznie material uwierzytelniajacy

**AuthLoginRequest** — Metodami dzialajacymi dzis sa HASLO i PIN; hello czeka na pochodzenie z domena po https. Login niesie sie przy metodzie password, bo konto ma nazwe ustalenie budowy nadana przy rejestracji; metody wlasciwe urzadzeniu tozsamosci nie potrzebuja, bo wskazuje ja material na urzadzeniu. Odmowa wraca kodem bledu not_authenticated. Proba nieudana NIE odmawia nastepnej — nakłada na nia ZWLOKE rosnaca wykladniczo (250 ms, 500, 1000, 2000, 4000, dalej rowno 5000 ms), zerowana przy pierwszym udanym wejsciu. Operator, ktory pomylil haslo trzy razy, wchodzi za czwartym; tylko czeka. Zadnego progu prob i zadnej odmowy 'za duzo prob' tu nie ma i nie bedzie — to byloby bramkowanie (rozstrzygniecie ustalenie budowy z 14.08.2026)

**AuthLoginRequest.Method** — Dzialajace dzis: password i pin. Metoda hello wymaga, by interfejs byl podany z pochodzenia z domena po https — pod http://127.0.0.1 jest niedostepna, bo WebAuthn wywodzi rp_id z pochodzenia dokumentu i zadna wartosc nie da sie podstawic z JavaScriptu

**AuthLoginRequest.Login** — Wymagana dla metody password; metody wlasciwe urzadzeniu (pin, hello) jej nie potrzebuja

**AuthLoginRequest.KeepSignedIn** — Nie jest bramka: znosi powtarzanie logowania, niczego nie blokuje

**AuthLoginResponse** — Metodami dzialajacymi dzis sa HASLO i PIN; hello czeka na pochodzenie z domena po https. Login niesie sie przy metodzie password, bo konto ma nazwe ustalenie budowy nadana przy rejestracji; metody wlasciwe urzadzeniu tozsamosci nie potrzebuja, bo wskazuje ja material na urzadzeniu. Odmowa wraca kodem bledu not_authenticated. Proba nieudana NIE odmawia nastepnej — nakłada na nia ZWLOKE rosnaca wykladniczo (250 ms, 500, 1000, 2000, 4000, dalej rowno 5000 ms), zerowana przy pierwszym udanym wejsciu. Operator, ktory pomylil haslo trzy razy, wchodzi za czwartym; tylko czeka. Zadnego progu prob i zadnej odmowy 'za duzo prob' tu nie ma i nie bedzie — to byloby bramkowanie (rozstrzygniecie ustalenie budowy z 14.08.2026)

**AuthMethodAddRequest** — Czynnosc USTAWIEN, wykonalna po zalogowaniu; hasla ta komenda nie zaklada, bo kotwica powstaje przy auth.register. Klucz Hello nie opuszcza TPM urzadzenia, wiec na kazdej maszynie zaklada sie go osobno — i tylko wtedy, gdy interfejs jest podawany z pochodzenia z domena po https. Dzis zalozyc mozna PIN

**AuthMethodAddRequest.Kind** — Zalozenie klucza hello wymaga pochodzenia z domena po https — pod http://127.0.0.1 rejestracja WebAuthn nie ma poprawnego rp_id

**AuthMethodAddResponse** — Czynnosc USTAWIEN, wykonalna po zalogowaniu; hasla ta komenda nie zaklada, bo kotwica powstaje przy auth.register. Klucz Hello nie opuszcza TPM urzadzenia, wiec na kazdej maszynie zaklada sie go osobno — i tylko wtedy, gdy interfejs jest podawany z pochodzenia z domena po https. Dzis zalozyc mozna PIN

**AuthMethodRemoveRequest** — Ani ostatniej metody, ani hasla zdjac sie nie da — haslo jest kotwica bramki

**AuthMethodRemoveResponse** — Ani ostatniej metody, ani hasla zdjac sie nie da — haslo jest kotwica bramki

**AuthTokenRefreshRequest** — W Danaco HUB dzieje sie to samo przy kazdym zadaniu HTTP; tu jest komenda, bo transportem jest jedno gniazdo WebSocket

**AuthTokenRefreshResponse** — W Danaco HUB dzieje sie to samo przy kazdym zadaniu HTTP; tu jest komenda, bo transportem jest jedno gniazdo WebSocket

**DeviceListRequest** — Konto jest jedno, urzadzen dowolnie wiele: komputery, telefony, tablety. Kazde niesie wlasny token dostepu, wiec wykaz jest miejscem, w ktorym Operator widzi, co ma dostep do platformy

**DeviceListResponse** — Konto jest jedno, urzadzen dowolnie wiele: komputery, telefony, tablety. Kazde niesie wlasny token dostepu, wiec wykaz jest miejscem, w ktorym Operator widzi, co ma dostep do platformy

**DeviceRevokeRequest** — Zmiana idzie do wszystkich polaczonych urzadzen zdarzeniem device.changed, wiec Operator widzi skutek natychmiast na pozostalych ekranach. Tej samej drogi uzywa odzyskanie konta, ktore uniewaznia tokeny wydane przed zmiana hasla

**DeviceRevokeResponse** — Zmiana idzie do wszystkich polaczonych urzadzen zdarzeniem device.changed, wiec Operator widzi skutek natychmiast na pozostalych ekranach. Tej samej drogi uzywa odzyskanie konta, ktore uniewaznia tokeny wydane przed zmiana hasla

**AppsDeploymentListRequest** — Bez tej komendy historia wdrozen znikala po odswiezeniu okna, bo klient odbudowywal ja wylacznie ze zdarzen biezacej sesji

**AppsDeploymentListResponse** — Bez tej komendy historia wdrozen znikala po odswiezeniu okna, bo klient odbudowywal ja wylacznie ze zdarzen biezacej sesji

**AppsArchitectureGetRequest** — Odpowiednik odczytu dla apps.architecture.define — bez niego okno po odswiezeniu nie wie, co Operator wczesniej zdefiniowal

**AppsArchitectureGetResponse** — Odpowiednik odczytu dla apps.architecture.define — bez niego okno po odswiezeniu nie wie, co Operator wczesniej zdefiniowal

**AppsWorkspaceListRequest** — Odpowiednik odczytu dla apps.workspace.update

**AppsWorkspaceListResponse** — Odpowiednik odczytu dla apps.workspace.update

**AppsWorkspaceListResponse.Files** — Kształt DeveloperFile — ten sam, którym apps.workspace.update oddaje plik po zapisie; druga struktura byłaby drugą prawdą o tym samym bycie

**DesignAssetUploadRequest** — Tresc wchodzi do magazynu rdzenia pod suma kontrolna — zasob przestaje zalezec od pliku, ktory Operator moze nadpisac albo skasowac. Uzyj, gdy zasob JUZ ISTNIEJE: Operator ma plik albo wskazuje odsylacz. Zasob majacy dopiero powstac zamawia sie komenda design.asset.generate

**DesignAssetUploadResponse** — Tresc wchodzi do magazynu rdzenia pod suma kontrolna — zasob przestaje zalezec od pliku, ktory Operator moze nadpisac albo skasowac. Uzyj, gdy zasob JUZ ISTNIEJE: Operator ma plik albo wskazuje odsylacz. Zasob majacy dopiero powstac zamawia sie komenda design.asset.generate

**DesignAssetRemoveRequest** — Rodzaj zmiany deleted zdarzenia design.asset.changed nie mial dotad zadnego nadawcy

**DesignAssetRemoveResponse** — Rodzaj zmiany deleted zdarzenia design.asset.changed nie mial dotad zadnego nadawcy

**ImageTransformRequest** — Wynikiem jest NOWY zasob w magazynie — zrodlo zostaje nietkniete, wiec model moze probowac bez ryzyka

**ImageTransformRequest.WindowId** — Model zna je z wlasnego zasiegu i podaje, zeby zasob byl WIDOCZNY dla Operatora w Assets Panel. Brak nie wstrzymuje czynnosci: bajty i tak trafiaja do magazynu pod suma kontrolna, ale zasob nie pojawi sie w wykazie okna

**ImageTransformResponse** — Wynikiem jest NOWY zasob w magazynie — zrodlo zostaje nietkniete, wiec model moze probowac bez ryzyka

**ImageAdjustRequest** — To jest RETUSZ, o ktory model prosi w rozmowie

**ImageAdjustRequest.WindowId** — Model zna je z wlasnego zasiegu i podaje, zeby zasob byl WIDOCZNY dla Operatora w Assets Panel. Brak nie wstrzymuje czynnosci: bajty i tak trafiaja do magazynu pod suma kontrolna, ale zasob nie pojawi sie w wykazie okna

**ImageAdjustResponse** — To jest RETUSZ, o ktory model prosi w rozmowie

**ImageConvertRequest** — Sluzy przygotowaniu zasobu do wydania — model uzywa jej przed osadzeniem grafiki na stronie

**ImageConvertRequest.WindowId** — Model zna je z wlasnego zasiegu i podaje, zeby zasob byl WIDOCZNY dla Operatora w Assets Panel. Brak nie wstrzymuje czynnosci: bajty i tak trafiaja do magazynu pod suma kontrolna, ale zasob nie pojawi sie w wykazie okna

**ImageConvertResponse** — Sluzy przygotowaniu zasobu do wydania — model uzywa jej przed osadzeniem grafiki na stronie

**MediaTranscodeRequest** — Wynikiem jest nowy zasob w magazynie

**MediaTranscodeRequest.WindowId** — Model zna je z wlasnego zasiegu i podaje, zeby zasob byl WIDOCZNY dla Operatora w Assets Panel. Brak nie wstrzymuje czynnosci: bajty i tak trafiaja do magazynu pod suma kontrolna, ale zasob nie pojawi sie w wykazie okna

**MediaTranscodeResponse** — Wynikiem jest nowy zasob w magazynie

**DocumentConvertRequest** — Model uzywa jej, gdy Operator prosi o dokument, a nie o tekst w oknie

**DocumentConvertRequest.WindowId** — Model zna je z wlasnego zasiegu i podaje, zeby zasob byl WIDOCZNY dla Operatora w Assets Panel. Brak nie wstrzymuje czynnosci: bajty i tak trafiaja do magazynu pod suma kontrolna, ale zasob nie pojawi sie w wykazie okna

**DocumentConvertResponse** — Model uzywa jej, gdy Operator prosi o dokument, a nie o tekst w oknie

**DocumentTextExtractRequest** — Model uzywa jej, zeby PRZECZYTAC to, co dostal jako plik

**DocumentTextExtractResponse** — Model uzywa jej, zeby PRZECZYTAC to, co dostal jako plik

**ArchivePackRequest** — Sluzy wydaniu pracy Operatorowi jednym plikiem

**ArchivePackRequest.WindowId** — Model zna je z wlasnego zasiegu i podaje, zeby zasob byl WIDOCZNY dla Operatora w Assets Panel. Brak nie wstrzymuje czynnosci: bajty i tak trafiaja do magazynu pod suma kontrolna, ale zasob nie pojawi sie w wykazie okna

**ArchivePackResponse** — Sluzy wydaniu pracy Operatorowi jednym plikiem

**ArchiveUnpackRequest.WindowId** — Model zna je z wlasnego zasiegu i podaje, zeby zasob byl WIDOCZNY dla Operatora w Assets Panel. Brak nie wstrzymuje czynnosci: bajty i tak trafiaja do magazynu pod suma kontrolna, ale zasob nie pojawi sie w wykazie okna

**MailAccountListRequest** — Bez skonfigurowanej skrzynki rodzina mail.* odmawia, nazywajac brak — rdzen poczty nie zmysla

**MailAccountListResponse** — Bez skonfigurowanej skrzynki rodzina mail.* odmawia, nazywajac brak — rdzen poczty nie zmysla

**MailMessageListRequest** — Uzyj do odnalezienia sprawy, o ktorej mowi Operator —na przykład listu od konkretnej osoby z dzisiaj

**MailMessageListResponse** — Uzyj do odnalezienia sprawy, o ktorej mowi Operator —na przykład listu od konkretnej osoby z dzisiaj

**MailMessageGetRequest** — Uzyj, gdy masz przeanalizowac list, a nie tylko go odnalezc. Zalaczniki trafiaja do magazynu rdzenia i wracaja jako zasoby, wiec model moze je dalej przetworzyc narzedziami obrazu i dokumentow

**MailMessageGetResponse** — Uzyj, gdy masz przeanalizowac list, a nie tylko go odnalezc. Zalaczniki trafiaja do magazynu rdzenia i wracaja jako zasoby, wiec model moze je dalej przetworzyc narzedziami obrazu i dokumentow

**MailDraftSaveRequest** — Szkic jest krokiem POSREDNIM: Operator widzi go w swojej poczcie, zanim cokolwiek wyjdzie w swiat

**MailDraftSaveResponse** — Szkic jest krokiem POSREDNIM: Operator widzi go w swojej poczcie, zanim cokolwiek wyjdzie w swiat

**MailSendRequest** — TO JEST JEDYNA KOMENDA RODZINY, KTOREJ SKUTEK WYCHODZI POZA MASZYNE OPERATORA i nie da sie go cofnac

**MailSendResponse** — TO JEST JEDYNA KOMENDA RODZINY, KTOREJ SKUTEK WYCHODZI POZA MASZYNE OPERATORA i nie da sie go cofnac

**MailAccountAddRequest** — Rdzen nie zaklada konta pocztowego i nie stawia serwera — bierze skrzynke, ktora juz istnieje: na urzadzeniu albo w chmurze. Poswiadczenie idzie do sejfu, nie do bazy

**MailAccountAddResponse** — Rdzen nie zaklada konta pocztowego i nie stawia serwera — bierze skrzynke, ktora juz istnieje: na urzadzeniu albo w chmurze. Poswiadczenie idzie do sejfu, nie do bazy

**MailAccountDiscoverRequest** — Niczego nie podpina i nie siega po haslo — oddaje to, co znalazl, zeby Operator nie przepisywal nastaw recznie

**MailAccountDiscoverResponse** — Niczego nie podpina i nie siega po haslo — oddaje to, co znalazl, zeby Operator nie przepisywal nastaw recznie

**MailAccountRemoveRequest** — Skrzynki u dostawcy nie tyka

**MailAccountRemoveResponse** — Skrzynki u dostawcy nie tyka

**ImageUpscaleRequest** — Zastepuje wyspecjalizowane modele powiekszajace; bez zainstalowanego silnika odmawia, nazywajac brak — nigdy nie oddaje zwyklego rozciagniecia jako powiekszenia

**ImageUpscaleResponse** — Zastepuje wyspecjalizowane modele powiekszajace; bez zainstalowanego silnika odmawia, nazywajac brak — nigdy nie oddaje zwyklego rozciagniecia jako powiekszenia

**ImageBackgroundRemoveRequest** — Zastepuje wyspecjalizowane narzedzia wycinania; bez silnika odmawia, nazywajac brak

**ImageBackgroundRemoveResponse** — Zastepuje wyspecjalizowane narzedzia wycinania; bez silnika odmawia, nazywajac brak

**KnowledgeIndexRequest** — Wyszukiwanie po slowach juz dziala (library.file.search); to jest droga do wyszukiwania po SENSIE

**KnowledgeIndexResponse** — Wyszukiwanie po slowach juz dziala (library.file.search); to jest droga do wyszukiwania po SENSIE

**KnowledgeSearchRequest.Rerank** — Podobienstwo wektorow jest pierwszym przebiegiem: pytanie i fragment licza sie osobno i spotykaja dopiero jako dwie liczby. Przesiew jest przebiegiem drugim — czyta pytanie RAZEM z fragmentem i uklada kandydatow na nowo. Kosztuje wczytanie drugiego modelu, wiec wchodzi na zadanie, a nie zawsze

**KnowledgeSearchRequest.RerankCandidates** — Liczba wieksza od `limit` jest tu sensem rzeczy: przesiew moze wyniesc na czolo fragment, ktory po samych wektorach byl dwudziesty. Bez `rerank` nie znaczy nic

**KnowledgeSearchResponse.Reranked** — Bez tego pola nie da sie odroznic odpowiedzi przesianej od odpowiedzi z samego pierwszego przebiegu, a `score` w obu ma ten sam ksztalt

**KnowledgeImageSearchRequest** — Rodzina knowledge.* prowadzi dotad wylacznie tekst; ta komenda jest jej osia obrazu i przeglada obrazy biblioteki Operatora

**KnowledgeImageSearchResponse** — Rodzina knowledge.* prowadzi dotad wylacznie tekst; ta komenda jest jej osia obrazu i przeglada obrazy biblioteki Operatora

**KnowledgeImageSearchResponse.Examined** — Wynik pusty przy zerze znaczy brak obrazow, a nie brak trafienia — dwie rozne odpowiedzi, ktorych bez tej liczby nie da sie rozroznic

**SessionToolAttachRequest** — Definicji eksperta NIE RUSZA (rozstrzygniecie ustalenie budowy, rozdz. 5): dolozenie zyje w stanie sesji, przezywa rozlaczenie klienta i konczy sie wraz z sesja albo z usunieciem rozmowy

**SessionToolAttachResponse** — Definicji eksperta NIE RUSZA (rozstrzygniecie ustalenie budowy, rozdz. 5): dolozenie zyje w stanie sesji, przezywa rozlaczenie klienta i konczy sie wraz z sesja albo z usunieciem rozmowy

**SessionToolDetachRequest** — Zestaw wraca do podstawy z definicji eksperta; sama definicja nie zmienia sie ani o joto, bo nigdy nie byla zmieniana

**SessionToolDetachResponse** — Zestaw wraca do podstawy z definicji eksperta; sama definicja nie zmienia sie ani o joto, bo nigdy nie byla zmieniana

**SessionToolListRequest** — Zestaw narzedzi tury to definicja eksperta PLUS te dolozenia — bez tego odczytu druga polowa zestawu bylaby niewidoczna

**SessionToolListResponse** — Zestaw narzedzi tury to definicja eksperta PLUS te dolozenia — bez tego odczytu druga polowa zestawu bylaby niewidoczna

**ToolsCatalogListRequest** — Wykaz liczy setki pozycji, wiec zadanie niesie zawezenie tekstem, rodzajem i grupa; dopasowanie idzie takze SRODKIEM nazwy, bo przy przedrostkach zrodla szukanie od poczatku byloby bezuzyteczne (rozdz. 7.2)

**ToolsCatalogListResponse** — Wykaz liczy setki pozycji, wiec zadanie niesie zawezenie tekstem, rodzajem i grupa; dopasowanie idzie takze SRODKIEM nazwy, bo przy przedrostkach zrodla szukanie od poczatku byloby bezuzyteczne (rozdz. 7.2)

**AgentConnectorListRequest** — Odpowiednik `agent.plugin.list` dla drugiego z dwoch bytow

**AgentConnectorListResponse** — Odpowiednik `agent.plugin.list` dla drugiego z dwoch bytow

**AgentVersionGetRequest** — `AgentVersion` niesie etykiete, opis zmiany, sprawce i sume kontrolna, ale NIE niesie tresci — bez tej komendy podgladu wersji nie ma z czego zlozyc, a przywrocenie jest jedynym sposobem zobaczenia, co w wersji stalo

**AgentVersionGetResponse** — `AgentVersion` niesie etykiete, opis zmiany, sprawce i sume kontrolna, ale NIE niesie tresci — bez tej komendy podgladu wersji nie ma z czego zlozyc, a przywrocenie jest jedynym sposobem zobaczenia, co w wersji stalo

**AgentAssignmentListRequest** — `workspace.agent.assign` zapisuje przynaleznosc, ale zadna komenda nie odczytuje jej OD STRONY EKSPERTA, wiec biblioteka nie ma skad wziac liczby

**AgentAssignmentListResponse** — `workspace.agent.assign` zapisuje przynaleznosc, ale zadna komenda nie odczytuje jej OD STRONY EKSPERTA, wiec biblioteka nie ma skad wziac liczby

**AgentModulesSetRequest** — LISTA PUSTA `[]` ZNACZY BRAK OGRANICZENIA, czyli dostepnosc wszedzie, i jest to stan wyjsciowy; pominiecie pola znaczy „bez zmiany". Ta sama zasada co przy `memoryLevels`, gdzie zbior pusty tez niesie znaczenie wlasne

**AgentModulesSetResponse** — LISTA PUSTA `[]` ZNACZY BRAK OGRANICZENIA, czyli dostepnosc wszedzie, i jest to stan wyjsciowy; pominiecie pola znaczy „bez zmiany". Ta sama zasada co przy `memoryLevels`, gdzie zbior pusty tez niesie znaczenie wlasne

**AgentIsolationSetRequest** — Wartosc obowiazuje wszedzie, gdzie ekspert dziala, dopoki nie nadpisze jej regula zapisana na poziomie zasiegu bardziej szczegolowym

**AgentIsolationSetResponse** — Wartosc obowiazuje wszedzie, gdzie ekspert dziala, dopoki nie nadpisze jej regula zapisana na poziomie zasiegu bardziej szczegolowym

**AgentPolicyGetRequest** — Transparentnosc, nie bramka — komenda niczego nie zapisuje i niczego nie rozstrzyga

**AgentPolicyGetResponse** — Transparentnosc, nie bramka — komenda niczego nie zapisuje i niczego nie rozstrzyga

**AgentSubagentSetRequest** — Granica pietnastu jest parametrem technicznym platformy, nie decyzja produktowa — ta sama, ktora zna `subagent.spawn`

**AgentSubagentSetResponse** — Granica pietnastu jest parametrem technicznym platformy, nie decyzja produktowa — ta sama, ktora zna `subagent.spawn`

**AgentPermissionRemoveRequest** — `agent.permission.set` wylacznie USTAWIA wartosc wpisu, wiec po pierwszym zawezeniu wiersz zakresu zostaje w wykazie na zawsze, choc jego wartosc wrocila do przyznanej

**AgentPermissionRemoveResponse** — `agent.permission.set` wylacznie USTAWIA wartosc wpisu, wiec po pierwszym zawezeniu wiersz zakresu zostaje w wykazie na zawsze, choc jego wartosc wrocila do przyznanej

**ChannelCheckRequest** — Narzedzie pomocnicze, nie bramka: wynik nie warunkuje zapisu eksperta ani wyslania tury

**ChannelCheckResponse** — Narzedzie pomocnicze, nie bramka: wynik nie warunkuje zapisu eksperta ani wyslania tury

**ChannelCredentialStatusRequest** — NIE ODDAJE TRESCI POSWIADCZENIA I ODDAWAC JEJ NIE MOZE

**ChannelCredentialStatusResponse** — NIE ODDAJE TRESCI POSWIADCZENIA I ODDAWAC JEJ NIE MOZE

**ExtensionSearchRequest** — Dopelnia extension.list, ktory zawezal wylacznie rodzajem i stanem zainstalowania; dopasowanie idzie takze SRODKIEM nazwy, bo pozycje niosa przedrostki zrodla

**ExtensionSearchResponse** — Dopelnia extension.list, ktory zawezal wylacznie rodzajem i stanem zainstalowania; dopasowanie idzie takze SRODKIEM nazwy, bo pozycje niosa przedrostki zrodla

**ExtensionDetailGetRequest** — Extension niesie sam naglowek pozycji, wiec karta szczegolow nie ma dzis z czego powstac

**ExtensionDetailGetResponse** — Extension niesie sam naglowek pozycji, wiec karta szczegolow nie ma dzis z czego powstac

**ExtensionCollectionApplyRequest** — Wynik jest bilansem, nie potwierdzeniem: pozycja odrzucona wraca z powodem, zamiast znikac po cichu

**ExtensionCollectionApplyResponse** — Wynik jest bilansem, nie potwierdzeniem: pozycja odrzucona wraca z powodem, zamiast znikac po cichu

**ExtensionRegistryListRequest** — Wyliczenie ExtensionOrigin ma dzis dwie wartosci, wiec rejestr organizacji nie ma jak sie w katalogu pokazac

**ExtensionRegistryListResponse** — Wyliczenie ExtensionOrigin ma dzis dwie wartosci, wiec rejestr organizacji nie ma jak sie w katalogu pokazac

**ExtensionUpdateCheckRequest** — Extension niesie wlasna wersje, ale nie wersje dostepna, wiec wskaznika aktualizacji nie ma dzis z czego zlozyc

**ExtensionUpdateCheckResponse** — Extension niesie wlasna wersje, ale nie wersje dostepna, wiec wskaznika aktualizacji nie ma dzis z czego zlozyc

**ExtensionVersionPinRequest** — Wersja jest dzis polem opisowym, wiec nie ma czego przypiac

**ExtensionVersionPinResponse** — Wersja jest dzis polem opisowym, wiec nie ma czego przypiac

**ExtensionVersionRollbackRequest** — Wykonuje sie od razu; ustawienie extension.rollback.confirm wlacza potwierdzenie, ktore nie warunkuje wykonania

**ExtensionVersionRollbackResponse** — Wykonuje sie od razu; ustawienie extension.rollback.confirm wlacza potwierdzenie, ktore nie warunkuje wykonania

**ExtensionBundleInstallRequest** — Wynik jest bilansem: pozycja odrzucona wraca z powodem

**ExtensionBundleInstallResponse** — Wynik jest bilansem: pozycja odrzucona wraca z powodem

**ExtensionHistoryListRequest** — Rejestr niesie stan biezacy, nie droge, ktora do niego doprowadzila

**ExtensionHistoryListResponse** — Rejestr niesie stan biezacy, nie droge, ktora do niego doprowadzila

**ExtensionAdminBulkRequest** — Wynik jest bilansem przyjetych i odrzuconych, nie pojedynczym potwierdzeniem

**ExtensionAdminBulkResponse** — Wynik jest bilansem przyjetych i odrzuconych, nie pojedynczym potwierdzeniem

**ExtensionToolListRequest** — Odpowiedz na tools/list, resources/list i prompts/list nie ma dzis w kontrakcie ani komendy, ani struktury, w ktora mialaby wejsc

**ExtensionToolListResponse** — Odpowiedz na tools/list, resources/list i prompts/list nie ma dzis w kontrakcie ani komendy, ani struktury, w ktora mialaby wejsc

**ExtensionToolCallRequest** — Wywolanie jest probne: idzie poza kontekstem eksperta i niczego mu nie przypisuje

**ExtensionToolCallResponse** — Wywolanie jest probne: idzie poza kontekstem eksperta i niczego mu nie przypisuje

**ExtensionTransportSetRequest** — Transport da sie dzis wpisac wylacznie do nieprzezroczystej konfiguracji, wiec rdzen i okno moga rozumiec to pole inaczej

**ExtensionTransportSetResponse** — Transport da sie dzis wpisac wylacznie do nieprzezroczystej konfiguracji, wiec rdzen i okno moga rozumiec to pole inaczej

**ExtensionCredentialBindRequest** — Moduł operuje na referencjach; tresc poswiadczenia nigdy nie opuszcza rdzenia i nie wchodzi w to zadanie

**ExtensionCredentialBindResponse** — Moduł operuje na referencjach; tresc poswiadczenia nigdy nie opuszcza rdzenia i nie wchodzi w to zadanie

**ExtensionSecretListRequest** — Oddaje wylacznie odwolania i metryke; tresc sekretu nie opuszcza rdzenia

**ExtensionSecretListResponse** — Oddaje wylacznie odwolania i metryke; tresc sekretu nie opuszcza rdzenia

**ExtensionPermissionListRequest** — Extension niesie nieprzezroczysta konfiguracje, ktora deklaracja uprawnien nie jest

**ExtensionPermissionListResponse** — Extension niesie nieprzezroczysta konfiguracje, ktora deklaracja uprawnien nie jest

**ExtensionPermissionGrantRequest** — Nadanie jest jedyna kontrola; brak nadania nie wstrzymuje instalacji ani wlaczenia (zasada zero blokad)

**ExtensionPermissionGrantResponse** — Nadanie jest jedyna kontrola; brak nadania nie wstrzymuje instalacji ani wlaczenia (zasada zero blokad)

**ExtensionManifestScanRequest** — Wynik jest SYGNALEM, nie brama: nie wstrzymuje instalacji ani wlaczenia

**ExtensionManifestScanResponse** — Wynik jest SYGNALEM, nie brama: nie wstrzymuje instalacji ani wlaczenia

**AppsStageListRequest** — Odpowiednik odczytu dla zdarzenia apps.build.changed — bez niego wykaz etapow zaczyna pusty przy kazdym wejsciu w modul, choc rdzen trzyma etapy w bazie

**AppsStageListResponse** — Odpowiednik odczytu dla zdarzenia apps.build.changed — bez niego wykaz etapow zaczyna pusty przy kazdym wejsciu w modul, choc rdzen trzyma etapy w bazie

**AppsStageSaveRequest** — AppStage.ownerAgentId jedzie dzis wylacznie w kierunku od rdzenia, wiec przypisania wykonawcy nie ma jak zlecic

**AppsStageSaveResponse** — AppStage.ownerAgentId jedzie dzis wylacznie w kierunku od rdzenia, wiec przypisania wykonawcy nie ma jak zlecic

**AppsProductLinkListRequest** — Powiazania nie sa czynne domyslnie — kazde jest swiadoma decyzja Operatora

**AppsProductLinkListResponse** — Powiazania nie sa czynne domyslnie — kazde jest swiadoma decyzja Operatora

**AppsArchitectureValidateRequest** — Dzis pole validationIssues jest zapisywane i oddawane bez wyliczenia, wiec jego pustka nie orzeka o poprawnosci ukladu. Zastrzezenia sa OSTRZEZENIEM, nie brama — nie wstrzymuja pracy w warsztatach

**AppsArchitectureValidateResponse** — Dzis pole validationIssues jest zapisywane i oddawane bez wyliczenia, wiec jego pustka nie orzeka o poprawnosci ukladu. Zastrzezenia sa OSTRZEZENIEM, nie brama — nie wstrzymuja pracy w warsztatach

**AppsArchitectureVersionListRequest** — AppArchitecture.version jest liczba, ale historii kontrakt nie niesie, wiec porownania nie ma z czym zestawic

**AppsArchitectureVersionListResponse** — AppArchitecture.version jest liczba, ale historii kontrakt nie niesie, wiec porownania nie ma z czym zestawic

**AppsPreviewStartRequest** — Podglad na zywo jest narzedziem warstwy 1 Frontend Workspace, a rdzen nie stawia dzis serwera produktu

**AppsPreviewStartResponse** — Podglad na zywo jest narzedziem warstwy 1 Frontend Workspace, a rdzen nie stawia dzis serwera produktu

**AppsEndpointListRequest** — AppComponent niesie apiContract jako JEDEN tekst, wiec punkt koncowy nie jest dzis bytem, ktory dalby sie wyliczyc

**AppsEndpointListResponse** — AppComponent niesie apiContract jako JEDEN tekst, wiec punkt koncowy nie jest dzis bytem, ktory dalby sie wyliczyc

**AppsEnvironmentListRequest** — Wyliczenie AppDeployEnvironment ma trzy wartosci stale, a opracowanie zada dowolnej liczby srodowisk konfigurowanych przez Operatora

**AppsEnvironmentListResponse** — Wyliczenie AppDeployEnvironment ma trzy wartosci stale, a opracowanie zada dowolnej liczby srodowisk konfigurowanych przez Operatora

**AppsEnvironmentVariableListRequest** — Wartosc sekretna oddawana jest jako referencja, nigdy jako tresc

**AppsEnvironmentVariableListResponse** — Wartosc sekretna oddawana jest jako referencja, nigdy jako tresc

**AppsServiceLogReadRequest** — Podglad na zywo idzie zdarzeniem apps.service.log

**AppsServiceLogReadResponse** — Podglad na zywo idzie zdarzeniem apps.service.log

**AppsArtifactListRequest** — AppDeployment niesie wersje i odnosnik logu, ale artefakt nie jest dzis bytem umowy

**AppsArtifactListResponse** — AppDeployment niesie wersje i odnosnik logu, ale artefakt nie jest dzis bytem umowy

**AppsDeploymentLogReadRequest** — AppDeployment niesie logRef — sam odnosnik, bez tresci — wiec podgladu logu nie ma dzis czym wypelnic. Podglad na zywo idzie zdarzeniem apps.deployment.log

**AppsDeploymentLogReadResponse** — AppDeployment niesie logRef — sam odnosnik, bez tresci — wiec podgladu logu nie ma dzis czym wypelnic. Podglad na zywo idzie zdarzeniem apps.deployment.log

**AppsPackageValidateRequest** — Zastrzezenia sa OSTRZEZENIEM, nie brama: nie wstrzymuja publikacji

**AppsPackageValidateResponse** — Zastrzezenia sa OSTRZEZENIEM, nie brama: nie wstrzymuja publikacji

**AppsPackageSignRequest** — Klucz wydawcy lezy w warstwie sekretow i nie wchodzi w to zadanie ani w odpowiedz

**AppsPackageSignResponse** — Klucz wydawcy lezy w warstwie sekretow i nie wchodzi w to zadanie ani w odpowiedz

**SpeechAudioUploadRequest** — Jest to brak platformowy, nie brak jednego okna: speech.transcribe bierze audioRef, czyli SCIEZKE PLIKU na maszynie silnika, a nagranie z mikrofonu istnieje wylacznie jako bajty w pamieci karty. Ten sam brak wstrzymuje dyktowanie w oknie komunikacji (client/src/okno-komunikacji/dyktowanie/dostarczenie-nagrania.ts) i mikrofon Voice Console — jedna komenda otwiera obie drogi naraz. Dzwiek nie opuszcza maszyny rdzenia: bajty ida do magazynu nagran, nie do sieci

**SpeechAudioUploadRequest.Retain** — Falsz kasuje je po rozpoznaniu; wartosc domyslna bierze ustawienie mowa_zapis_nagran

**SpeechAudioUploadResponse** — Jest to brak platformowy, nie brak jednego okna: speech.transcribe bierze audioRef, czyli SCIEZKE PLIKU na maszynie silnika, a nagranie z mikrofonu istnieje wylacznie jako bajty w pamieci karty. Ten sam brak wstrzymuje dyktowanie w oknie komunikacji (client/src/okno-komunikacji/dyktowanie/dostarczenie-nagrania.ts) i mikrofon Voice Console — jedna komenda otwiera obie drogi naraz. Dzwiek nie opuszcza maszyny rdzenia: bajty ida do magazynu nagran, nie do sieci

**SpeechAudioFetchRequest** — Domyka pare do speech.audio.upload i do translate.speech.synthesize: dzis rdzen oddaje odnosniki (AssistantActivityEntry.audioRef, AssistantVoiceCommandResponse.speechRef), lecz nie ma komendy, ktora by po nie siegnela — odslucha nie ma czym wykonac

**SpeechAudioFetchResponse** — Domyka pare do speech.audio.upload i do translate.speech.synthesize: dzis rdzen oddaje odnosniki (AssistantActivityEntry.audioRef, AssistantVoiceCommandResponse.speechRef), lecz nie ma komendy, ktora by po nie siegnela — odslucha nie ma czym wykonac

**SpeechWakeGetRequest** — Brak modelu frazy jest ODPOWIEDZIA, nie awaria — tak samo jak brak silnika w speech.availability.get

**SpeechWakeGetResponse** — Brak modelu frazy jest ODPOWIEDZIA, nie awaria — tak samo jak brak silnika w speech.availability.get

**SpeechWakeSetRequest** — Pola pominiete zostaja bez zmian

**SpeechWakeSetResponse** — Pola pominiete zostaja bez zmian

**SpeechListenStartRequest** — Rdzen slucha strumienia i oglasza zdarzeniami: czesciowa transkrypcje oraz wykrycie frazy wybudzajacej. Nasluch nie jest bramka — Operator zatrzymuje go speech.listen.stop, a rdzen nie zatrzymuje go sam

**SpeechListenStartResponse** — Rdzen slucha strumienia i oglasza zdarzeniami: czesciowa transkrypcje oraz wykrycie frazy wybudzajacej. Nasluch nie jest bramka — Operator zatrzymuje go speech.listen.stop, a rdzen nie zatrzymuje go sam

**SpeechListenStopRequest** — Zatrzymanie nasluchu, ktorego nie ma, nie jest bledem

**SpeechListenStopResponse** — Zatrzymanie nasluchu, ktorego nie ma, nie jest bledem

**AssistantActivityFlagRequest** — Dzis wyroznienie wpisu nie mialoby gdzie zamieszkac: AssistantActivityEntry nie niesie takiego pola, wiec okno pokazywaloby wyroznienie, ktorego rdzen nie pamieta

**AssistantActivityFlagResponse** — Dzis wyroznienie wpisu nie mialoby gdzie zamieszkac: AssistantActivityEntry nie niesie takiego pola, wiec okno pokazywaloby wyroznienie, ktorego rdzen nie pamieta

**MemoryRetentionGetRequest** — Odczyt do pary z memory.retention.set

**MemoryRetentionGetResponse** — Odczyt do pary z memory.retention.set

**MemoryRetentionSetRequest** — Zasada obejmuje zapisy kolejne, a wpisow zastanych nie rusza wstecz — inaczej zmiana reguly kasowalaby ustalenia, na ktore Operator sie nie umawial

**MemoryRetentionSetResponse** — Zasada obejmuje zapisy kolejne, a wpisow zastanych nie rusza wstecz — inaczej zmiana reguly kasowalaby ustalenia, na ktore Operator sie nie umawial

**MemoryContextListRequest** — Dzis memory.toggle przestawia poziomy zasiegu i to jedyny kontekst przelaczalny; nazwy zestawu nie ma gdzie odlozyc

**MemoryContextListResponse** — Dzis memory.toggle przestawia poziomy zasiegu i to jedyny kontekst przelaczalny; nazwy zestawu nie ma gdzie odlozyc

**MemoryContextSaveRequest** — Puste contextId zaklada nowy, podane zmienia istniejacy — tak samo jak memory.set

**MemoryContextSaveResponse** — Puste contextId zaklada nowy, podane zmienia istniejacy — tak samo jak memory.set

**MemoryContextActivateRequest** — Aktywacja nie kasuje kontekstu poprzedniego — wraca sie do niego tym samym wywolaniem

**MemoryContextActivateResponse** — Aktywacja nie kasuje kontekstu poprzedniego — wraca sie do niego tym samym wywolaniem

**MemoryContextDeleteRequest** — Wpisow pamieci nie kasuje — kontekst jest zestawem wskazan, a nie ustalenie budowy tresci

**MemoryContextDeleteResponse** — Wpisow pamieci nie kasuje — kontekst jest zestawem wskazan, a nie ustalenie budowy tresci

**ToolsScopeListRequest** — Dzis takiego zakresu nie ma dla profilu: agent.connector.* i agent.permission.set dotycza eksperta modulu Agents i nie siegaja profilu asystenta

**ToolsScopeListResponse** — Dzis takiego zakresu nie ma dla profilu: agent.connector.* i agent.permission.set dotycza eksperta modulu Agents i nie siegaja profilu asystenta

**ToolsScopeSetRequest** — Zakres jest nastawa zasiegu, nie bramka: stanem wyjsciowym jest pelny dostep bez limitu, zgodnie z zasada zero blokad

**ToolsScopeSetResponse** — Zakres jest nastawa zasiegu, nie bramka: stanem wyjsciowym jest pelny dostep bez limitu, zgodnie z zasada zero blokad

**ClipboardListRequest** — Schowek jest bytem klienta, ktorego rdzen dzis nie zna, wiec historia ginie wraz z karta — komenda daje jej trwalosc

**ClipboardListResponse** — Schowek jest bytem klienta, ktorego rdzen dzis nie zna, wiec historia ginie wraz z karta — komenda daje jej trwalosc

**ClipboardPushRequest** — Powtorzenie tresci identycznej nie mnozy wpisow — podnosi wpis zastany na czolo wykazu

**ClipboardPushResponse** — Powtorzenie tresci identycznej nie mnozy wpisow — podnosi wpis zastany na czolo wykazu

**ClipboardPinRequest** — Przypiety nie wygasa wraz z zasada retencji historii

**ClipboardPinResponse** — Przypiety nie wygasa wraz z zasada retencji historii

**SnippetListRequest** — Skrot dziala we wszystkich polach tekstowych platformy, wiec slownik nalezy do rdzenia, nie do jednego okna

**SnippetListResponse** — Skrot dziala we wszystkich polach tekstowych platformy, wiec slownik nalezy do rdzenia, nie do jednego okna

**SnippetSetRequest** — Puste snippetId zaklada nowy, podane zmienia istniejacy

**SnippetSetResponse** — Puste snippetId zaklada nowy, podane zmienia istniejacy

**LauncherHotkeyGetRequest** — Skrot globalny rejestruje powloka programu okiennego, wiec brak wsparcia powloki jest ODPOWIEDZIA, nie awaria — klient nie obiecuje wtedy skrotu, ktory nikogo nie obudzi

**LauncherHotkeyGetResponse** — Skrot globalny rejestruje powloka programu okiennego, wiec brak wsparcia powloki jest ODPOWIEDZIA, nie awaria — klient nie obiecuje wtedy skrotu, ktory nikogo nie obudzi

**LauncherHotkeyGetResponse.Hotkey** — Ctrl+Shift+Space; pusty znaczy brak nastawy

**LauncherHotkeySetRequest** — Skrot zajety przez inny program nie jest bledem zapisu — zapis zostaje, a odpowiedz mowi, ze rejestracja sie nie udala

**LauncherHotkeySetResponse** — Skrot zajety przez inny program nie jest bledem zapisu — zapis zostaje, a odpowiedz mowi, ze rejestracja sie nie udala

**ContextUsageGetRequest** — ZALEZNOSC: wymaga tokenizatora jako podsystemu rdzenia — bez niego liczba tokenow bylaby wartoscia wzieta znikad. Do czasu jego zbudowania komenda ma zwracac available falsz wraz z powodem, tak jak speech.availability.get przy braku silnika

**ContextUsageGetResponse** — ZALEZNOSC: wymaga tokenizatora jako podsystemu rdzenia — bez niego liczba tokenow bylaby wartoscia wzieta znikad. Do czasu jego zbudowania komenda ma zwracac available falsz wraz z powodem, tak jak speech.availability.get przy braku silnika

**DesignAssetContentGetRequest** — Pole uri zasobu jest sciezka w systemie plikow rdzenia (magazyn oddaje filepath), wiec przegladarka nie wczyta spod niego niczego — take wtedy, gdy zasob powstal bez zarzutu. Komenda dotyczy KAZDEGO zasobu magazynu, nie tylko obrazu: jeden magazyn obsluguje rodziny design.*, document.*, media.* i archive.*. Rdzen odmawia zasobu, ktorego nie zna, i zasobu, ktorego tresci nie ma pod suma kontrolna — nigdy nie oddaje bajtow zastepczych

**DesignAssetContentGetResponse** — Pole uri zasobu jest sciezka w systemie plikow rdzenia (magazyn oddaje filepath), wiec przegladarka nie wczyta spod niego niczego — take wtedy, gdy zasob powstal bez zarzutu. Komenda dotyczy KAZDEGO zasobu magazynu, nie tylko obrazu: jeden magazyn obsluguje rodziny design.*, document.*, media.* i archive.*. Rdzen odmawia zasobu, ktorego nie zna, i zasobu, ktorego tresci nie ma pod suma kontrolna — nigdy nie oddaje bajtow zastepczych

**DesignAssetExportRequest** — Rozni sie od image.convert celem: konwersja zaklada NOWY zasob w magazynie i tam sie konczy, eksport oddaje bajty gotowe do zapisania poza produktem

**DesignAssetExportResponse** — Rozni sie od image.convert celem: konwersja zaklada NOWY zasob w magazynie i tam sie konczy, eksport oddaje bajty gotowe do zapisania poza produktem

**DesignAssetExportBatchRequest** — Pokrywa eksport zbiorczy Assets Panel oraz zestawy rozmiarow kampanii; odmowa jednego zasobu NIE wstrzymuje pozostalych — wynik niesie bilans przyjetych i odrzuconych

**DesignAssetExportBatchResponse** — Pokrywa eksport zbiorczy Assets Panel oraz zestawy rozmiarow kampanii; odmowa jednego zasobu NIE wstrzymuje pozostalych — wynik niesie bilans przyjetych i odrzuconych

**DesignCollectionCreateRequest** — Etykiety zasobu juz sa (design.asset.tag.set), ale kolekcja jest bytem osobnym: ma nazwe, opis i porzadek, a etykieta jest tylko slowem

**DesignCollectionCreateResponse** — Etykiety zasobu juz sa (design.asset.tag.set), ale kolekcja jest bytem osobnym: ma nazwe, opis i porzadek, a etykieta jest tylko slowem

**DesignCollectionAssignRequest** — Zestaw jest DOKLADKA albo ODJECIEM, nie zastapieniem — inaczej niz przy etykietach, bo kolekcja bywa duza i przepisywanie jej w calosci przy kazdej zmianie jest droga do zgubienia zawartosci

**DesignCollectionAssignResponse** — Zestaw jest DOKLADKA albo ODJECIEM, nie zastapieniem — inaczej niz przy etykietach, bo kolekcja bywa duza i przepisywanie jej w calosci przy kazdej zmianie jest droga do zgubienia zawartosci

**DesignCollectionListRequest** — Bez tego kolekcja zalozona nie mialaby jak wrocic na ekran po odswiezeniu

**DesignCollectionListResponse** — Bez tego kolekcja zalozona nie mialaby jak wrocic na ekran po odswiezeniu

**DesignPromptTemplateSaveRequest** — Dzis historia promptow i szablony zyja w oknie do zamkniecia karty przegladarki i tyle o nich wiadomo

**DesignPromptTemplateSaveResponse** — Dzis historia promptow i szablony zyja w oknie do zamkniecia karty przegladarki i tyle o nich wiadomo

**DesignPromptHistoryListRequest** — Domyka tez prowenancje: dzis zasob niesie pole promptId, ktorego rdzen NIE wypelnia, bo nie ma przekladu klucza wiersza promptu na kod kontraktu — ta komenda ten przeklad wnosi

**DesignPromptHistoryListResponse** — Domyka tez prowenancje: dzis zasob niesie pole promptId, ktorego rdzen NIE wypelnia, bo nie ma przekladu klucza wiersza promptu na kod kontraktu — ta komenda ten przeklad wnosi

**DesignBoardVersionRestoreRequest** — Przywrocenie ZAKLADA nowa wersje z ukladu sprzed przywrocenia, zeby cofniecie sie samo nie kasowalo stanu, ktory Operator wlasnie porzucil

**DesignBoardVersionRestoreResponse** — Przywrocenie ZAKLADA nowa wersje z ukladu sprzed przywrocenia, zeby cofniecie sie samo nie kasowalo stanu, ktory Operator wlasnie porzucil

**DesignBoardExportRequest** — Dzis kompozycja jezdzi do rdzenia i z powrotem jako uklad warstw i nie ma drogi wyjscia poza rdzen

**DesignBoardExportResponse** — Dzis kompozycja jezdzi do rdzenia i z powrotem jako uklad warstw i nie ma drogi wyjscia poza rdzen

**DesignAnnotationSetRequest** — Pole note warstwy niesie JEDNO zdanie bez autora i bez watku; opracowanie wymaga watkow i oznaczen osob, a tego jedno pole nie unosi

**DesignAnnotationSetResponse** — Pole note warstwy niesie JEDNO zdanie bez autora i bez watku; opracowanie wymaga watkow i oznaczen osob, a tego jedno pole nie unosi

**DesignPresenceReportRequest** — Rdzen rozglasza je pozostalym zdarzeniem design.board.presence. Zgloszenie jest ULOTNE — nie zapisuje sie w bazie, bo polozenie kursora sprzed godziny nie jest wiedza o niczym

**DesignPresenceReportResponse** — Rdzen rozglasza je pozostalym zdarzeniem design.board.presence. Zgloszenie jest ULOTNE — nie zapisuje sie w bazie, bo polozenie kursora sprzed godziny nie jest wiedza o niczym

**DesignTokensetSaveRequest** — Dzis Tokens & System Panel czyta zetony z motywu obowiazujacego i nie ma ich gdzie odlozyc — kontrakt nie zna bytu zestawu zetonow

**DesignTokensetSaveResponse** — Dzis Tokens & System Panel czyta zetony z motywu obowiazujacego i nie ma ich gdzie odlozyc — kontrakt nie zna bytu zestawu zetonow

**DesignTokensetExportRequest** — Klient sklada dzis zmienne CSS, SCSS, konfiguracje Tailwind i modul JavaScript sam i nie potrzebuje do tego rdzenia; ta komenda jest potrzebna dla postaci, ktorych przegladarka zlozyc nie moze, oraz dla wydania idacego DO INNEGO MODULU zamiast do pliku

**DesignTokensetExportResponse** — Klient sklada dzis zmienne CSS, SCSS, konfiguracje Tailwind i modul JavaScript sam i nie potrzebuje do tego rdzenia; ta komenda jest potrzebna dla postaci, ktorych przegladarka zlozyc nie moze, oraz dla wydania idacego DO INNEGO MODULU zamiast do pliku

**DesignTokensetImportRequest** — Rdzen NIE nadpisuje motywu produktu — motyw jest wlasnoscia powloki; import zaklada byt obok niego i oddaje roznice wobec zetonow wskazanego motywu

**DesignTokensetImportResponse** — Rdzen NIE nadpisuje motywu produktu — motyw jest wlasnoscia powloki; import zaklada byt obok niego i oddaje roznice wobec zetonow wskazanego motywu

**DesignStyleguidePublishRequest** — Przegladarka sklada przewodnik sama i oddaje go plikiem, ale wydania go do Library albo Studio nie ma czym zlecic

**DesignStyleguidePublishResponse** — Przegladarka sklada przewodnik sama i oddaje go plikiem, ale wydania go do Library albo Studio nie ma czym zlecic

**ImageVectorizeRequest** — Rodzina image.* zna przeksztalcenie geometryczne, poprawke, konwersje formatu, powiekszenie i wyciecie tla — zamiany rastra na sciezki nie zna zadna z nich, a jest to funkcja opracowania modulu

**ImageVectorizeResponse** — Rodzina image.* zna przeksztalcenie geometryczne, poprawke, konwersje formatu, powiekszenie i wyciecie tla — zamiany rastra na sciezki nie zna zadna z nich, a jest to funkcja opracowania modulu

**ImageLayersSplitRequest** — Zasila rozdzielenie generacji na warstwy oraz zaznaczanie obiektu; bez silnika segmentacji ODMAWIA, nazywajac brak — nigdy nie oddaje calego obrazu jako jednej warstwy udajac rozklad

**ImageLayersSplitResponse** — Zasila rozdzielenie generacji na warstwy oraz zaznaczanie obiektu; bez silnika segmentacji ODMAWIA, nazywajac brak — nigdy nie oddaje calego obrazu jako jednej warstwy udajac rozklad

**ImageComposeRequest** — Pokrywa znak wodny, branding wsadowy, osadzenie w ramce urzadzenia i warstwy rastrowe — cztery funkcje opracowania, ktore wszystkie potrzebuja jednego: zlozenia obrazu na obrazie

**ImageComposeResponse** — Pokrywa znak wodny, branding wsadowy, osadzenie w ramce urzadzenia i warstwy rastrowe — cztery funkcje opracowania, ktore wszystkie potrzebuja jednego: zlozenia obrazu na obrazie

**RoundtableDebateGetRequest** — Obszar nie ma dzis zadnej komendy odczytu, wiec okno otwarte w trakcie debaty zna wylacznie to, co uslyszalo zdarzeniem od swojego otwarcia

**RoundtableDebateGetResponse** — Obszar nie ma dzis zadnej komendy odczytu, wiec okno otwarte w trakcie debaty zna wylacznie to, co uslyszalo zdarzeniem od swojego otwarcia

**StudioDocumentFormatSetRequest** — Dzis format nadaje rdzen przy wczytaniu i zadna komenda go nie zmienia: studio.document.save przyjmuje documentId, content, title i createVersion, ale nie format

**StudioDocumentFormatSetResponse** — Dzis format nadaje rdzen przy wczytaniu i zadna komenda go nie zmienia: studio.document.save przyjmuje documentId, content, title i createVersion, ale nie format

**StudioCommentAddRequest.Author** — Model zaklada komentarz podpisany jako model i ta droga, nie druga

**StudioCommentResolveRequest** — Watek nie znika: rozwiazanie jest stanem, nie usunieciem

**StudioCommentResolveResponse** — Watek nie znika: rozwiazanie jest stanem, nie usunieciem

**StudioTrackingSetRequest** — Przy wlaczonym sledzeniu kazdy zapis odklada wstawienia i usuniecia jako zmiany do decyzji, zamiast nadpisywac tresc

**StudioTrackingSetResponse** — Przy wlaczonym sledzeniu kazdy zapis odklada wstawienia i usuniecia jako zmiany do decyzji, zamiast nadpisywac tresc

**StudioTrackingDecideRequest** — Decyzja zapisuje sie w tresci dokumentu i zaklada wersje

**StudioTrackingDecideResponse** — Decyzja zapisuje sie w tresci dokumentu i zaklada wersje

**StudioAssetEmbedRequest** — Realizuje powiazanie Design do Studio

**StudioAssetEmbedResponse** — Realizuje powiazanie Design do Studio

**StudioOperationDeleteRequest** — Operacji fabrycznej nie usuwa — na nia odpowiada odmowa nazywajaca ten fakt

**StudioOperationDeleteResponse** — Operacji fabrycznej nie usuwa — na nia odpowiada odmowa nazywajaca ten fakt

**StudioChainRunRequest** — Przebieg prowadzi petla wykonawcza okna, a kazdy krok odklada wlasna propozycje zmiany

**StudioChainRunResponse** — Przebieg prowadzi petla wykonawcza okna, a kazdy krok odklada wlasna propozycje zmiany

**StudioBatchRunRequest** — Kazdy dokument dostaje wlasne zadanie petli; odmowa jednego nie wstrzymuje pozostalych

**StudioBatchRunResponse** — Kazdy dokument dostaje wlasne zadanie petli; odmowa jednego nie wstrzymuje pozostalych

**StudioAnnotationAddRequest** — Dzis czynnosc ta idzie droga generyczna window.action i wraca odmowa not_found, bo katalog akcji nie ma jej wiersza

**StudioAnnotationAddResponse** — Dzis czynnosc ta idzie droga generyczna window.action i wraca odmowa not_found, bo katalog akcji nie ma jej wiersza

**StudioDiffReportExportRequest** — Wynik jest zasobem magazynu rdzenia

**StudioDiffReportExportResponse** — Wynik jest zasobem magazynu rdzenia

**StudioDiffVisualRequest** — Sluzy tam, gdzie roznica tekstowa nie widzi zmiany ukladu

**StudioDiffVisualResponse** — Sluzy tam, gdzie roznica tekstowa nie widzi zmiany ukladu

**StudioSearchSemanticRequest** — Uzupelnia wyszukiwanie wzorca w studio.diff.compare

**StudioSearchSemanticResponse** — Uzupelnia wyszukiwanie wzorca w studio.diff.compare

**StudioDiffSourceRequest** — Sluzy kontroli, czy redakcja nie odeszla od zrodla

**StudioDiffSourceResponse** — Sluzy kontroli, czy redakcja nie odeszla od zrodla

**StudioProposalDecideRequest** — Dzis decyzja zapada wylacznie w kliencie i rdzen o niej nie wie, wiec propozycja zostaje w nim nierozstrzygnieta

**StudioProposalDecideResponse** — Dzis decyzja zapada wylacznie w kliencie i rdzen o niej nie wie, wiec propozycja zostaje w nim nierozstrzygnieta

**StudioBranchCreateRequest** — Sluzy prowadzeniu dwoch redakcji obok siebie

**StudioBranchCreateResponse** — Sluzy prowadzeniu dwoch redakcji obok siebie

**StudioBranchMergeRequest** — Konflikt nierozstrzygniety wraca w wyniku zamiast byc rozstrzygniety domyslem

**StudioBranchMergeResponse** — Konflikt nierozstrzygniety wraca w wyniku zamiast byc rozstrzygniety domyslem

**StudioPreviewRenderRequest** — Zwraca strony jako zasoby, wiec podglad pokazuje UKLAD, a nie sam tekst

**StudioPreviewRenderResponse** — Zwraca strony jako zasoby, wiec podglad pokazuje UKLAD, a nie sam tekst

**StudioPdfSplitRequest** — Kazda czesc jest osobnym zasobem

**StudioPdfSplitResponse** — Kazda czesc jest osobnym zasobem

**StudioSecurityEncryptRequest** — Czynnosc jest jawnym, odwracalnym ustawieniem Operatora i niczego nie warunkuje

**StudioSecurityEncryptResponse** — Czynnosc jest jawnym, odwracalnym ustawieniem Operatora i niczego nie warunkuje

**StudioSecurityRedactRequest** — Czynnosc jest nieodwracalna dla wyniku, dlatego zrodlo zostaje nietkniete

**StudioSecurityRedactResponse** — Czynnosc jest nieodwracalna dla wyniku, dlatego zrodlo zostaje nietkniete

**StudioSecuritySensitiveDetectRequest** — Czynnosc CZYTA i niczego nie zmienia

**StudioSecuritySensitiveDetectResponse** — Czynnosc CZYTA i niczego nie zmienia

**StudioIngestQueueAddRequest** — Dzis kolejka jest wylacznie kliencka i ginie z odswiezeniem okna

**StudioIngestQueueAddResponse** — Dzis kolejka jest wylacznie kliencka i ginie z odswiezeniem okna

**StudioIngestRecognizeRequest** — Komenda document.text.extract robi to samo waskim wejsciem — jeden jezyk, bez silnika, bez progu i bez ukladu

**StudioIngestRecognizeResponse** — Komenda document.text.extract robi to samo waskim wejsciem — jeden jezyk, bez silnika, bez progu i bez ukladu

**StudioIngestUrlRequest** — Migawke strony oddaje browser.snapshot.get, ale wymaga okna modulu Browser i nie prowadzi do dokumentu Studia

**StudioIngestUrlResponse** — Migawke strony oddaje browser.snapshot.get, ale wymaga okna modulu Browser i nie prowadzi do dokumentu Studia

**StudioIngestDeviceListRequest** — Bez tego wykazu wybor urzadzenia nie ma z czego powstac

**StudioIngestDeviceListResponse** — Bez tego wykazu wybor urzadzenia nie ma z czego powstac

**StudioStyleSaveRequest** — Zmiana stylu przestawia wszystkie miejsca dokumentu, ktore go uzywaja

**StudioStyleSaveResponse** — Zmiana stylu przestawia wszystkie miejsca dokumentu, ktore go uzywaja

**StudioStyleDeleteRequest** — Stylu fabrycznego nie usuwa — odpowiada odmowa nazywajaca powod

**StudioStyleDeleteResponse** — Stylu fabrycznego nie usuwa — odpowiada odmowa nazywajaca powod

**StudioPageSetupSetRequest** — Zmiana formatu przelicza uklad i oddaje bilans tego, co sie nie zmiescilo

**StudioPageSetupSetResponse** — Zmiana formatu przelicza uklad i oddaje bilans tego, co sie nie zmiescilo

**StudioTableStructureEditRequest** — Szerokosci kolumn zostaja policzone, nie zerowe

**StudioTableStructureEditResponse** — Szerokosci kolumn zostaja policzone, nie zerowe

**StudioDocumentImportPdfRequest** — Odzyskanie jest odtworzeniem, nie odczytem, wiec odpowiedz niesie bilans; PDF ze samych skanow kieruje na rozpoznanie tekstu

**StudioDocumentImportPdfResponse** — Odzyskanie jest odtworzeniem, nie odczytem, wiec odpowiedz niesie bilans; PDF ze samych skanow kieruje na rozpoznanie tekstu

**StudioDocumentCopyRequest** — Kopia jest osobnym dokumentem, nie drugim odwolaniem do tego samego

**StudioDocumentCopyResponse** — Kopia jest osobnym dokumentem, nie drugim odwolaniem do tego samego

**StudioDocumentExportFormatRequest** — Format ubozszy niz dokument oddaje wykaz cech pominietych

**StudioDocumentExportFormatResponse** — Format ubozszy niz dokument oddaje wykaz cech pominietych

**StudioTemplateDeleteRequest** — Szablonu fabrycznego nie usuwa — odpowiada odmowa nazywajaca powod

**StudioTemplateDeleteResponse** — Szablonu fabrycznego nie usuwa — odpowiada odmowa nazywajaca powod

**StudioLockAddRequest** — Blokada obowiazuje w rdzeniu, przed dotknieciem tresci

**StudioLockAddResponse** — Blokada obowiazuje w rdzeniu, przed dotknieciem tresci

**StudioLockRemoveRequest** — Blokade zdejmuje WYLACZNIE Operator — czynnosc modelu wraca odmowa nazywajaca powod

**StudioLockRemoveResponse** — Blokade zdejmuje WYLACZNIE Operator — czynnosc modelu wraca odmowa nazywajaca powod

**StudioJournalRevertRequest** — Czynnosc bedaca podstawa pozniejszej odmawia i nazywa zaleznosc, zamiast zostawic dokument w stanie niespojnym

**StudioJournalRevertResponse** — Czynnosc bedaca podstawa pozniejszej odmawia i nazywa zaleznosc, zamiast zostawic dokument w stanie niespojnym

**StudioModelChangesRevertRequest** — Nie jest to przywrocenie wersji sprzed, bo to skasowaloby prace Operatora

**StudioModelChangesRevertResponse** — Nie jest to przywrocenie wersji sprzed, bo to skasowaloby prace Operatora

**StudioAutosaveRunRequest** — Nieudany zapis wraca nazwany, nie przemilczany

**StudioAutosaveRunResponse** — Nieudany zapis wraca nazwany, nie przemilczany

**StudioViewSetRequest** — Gdzie da sie zrobic dwojako i obie drogi maja sens, wybor nalezy do Operatora i jest jawnym, odwracalnym ustawieniem

**StudioViewSetResponse** — Gdzie da sie zrobic dwojako i obie drogi maja sens, wybor nalezy do Operatora i jest jawnym, odwracalnym ustawieniem

**StudioMarkupAddRequest** — Znakowanie modelu jest podpisane jako model

**StudioMarkupAddResponse** — Znakowanie modelu jest podpisane jako model

**StudioMarkupTypeDeleteRequest** — Rodzaju fabrycznego nie usuwa — odpowiada odmowa nazywajaca powod

**StudioMarkupTypeDeleteResponse** — Rodzaju fabrycznego nie usuwa — odpowiada odmowa nazywajaca powod

**StudioAgentsSettingsGetRequest** — Oba narzedzia sa domyslnie wylaczone

**StudioAgentsSettingsGetResponse** — Oba narzedzia sa domyslnie wylaczone

**StudioAgentsSettingsSetRequest** — Wartosci ida zasiegami rodziny config, nie osobnym magazynem

**StudioAgentsSettingsSetResponse** — Wartosci ida zasiegami rodziny config, nie osobnym magazynem

**StudioPlanCreateRequest** — Sam rozklad niczego nie uruchamia

**StudioPlanCreateResponse** — Sam rozklad niczego nie uruchamia

**StudioPlanRunRequest** — Gdy nastawa petli jest wylaczona, wraca odmowa nazywajaca brak nastawy, a nie cisza

**StudioPlanRunResponse** — Gdy nastawa petli jest wylaczona, wraca odmowa nazywajaca brak nastawy, a nie cisza

**StudioAgentsClaimRequest** — Fragment zajety przez kogo innego wraca odmowa nazywajaca wykonawce i zakres

**StudioAgentsClaimResponse** — Fragment zajety przez kogo innego wraca odmowa nazywajaca wykonawce i zakres

**StudioAgentsConflictsListRequest** — Odlozone brzmienie nie przepada

**StudioAgentsConflictsListResponse** — Odlozone brzmienie nie przepada

**TerminalSessionCloseRequest** — Dzis zamkniecie karty zyje wylacznie w widoku klienta: powloka i jej procesy biegna dalej, a rdzen o zamknieciu nie wie

**TerminalSessionCloseResponse** — Dzis zamkniecie karty zyje wylacznie w widoku klienta: powloka i jej procesy biegna dalej, a rdzen o zamknieciu nie wie

**TerminalSessionListRequest** — Rdzen odtwarza karty przy starcie, ale klient po ponownym polaczeniu nie ma jak ich zobaczyc i zaczyna wykaz od pustego

**TerminalSessionListResponse** — Rdzen odtwarza karty przy starcie, ale klient po ponownym polaczeniu nie ma jak ich zobaczyc i zaczyna wykaz od pustego

**TerminalFileReadRequest** — Bez tej komendy klient czyta manifest projektu poleceniem powloki, wiec wykrycie zadan zalezy od programu wypisujacego plik i od skladni kazdej z szesciu powlok

**TerminalFileReadResponse** — Bez tej komendy klient czyta manifest projektu poleceniem powloki, wiec wykrycie zadan zalezy od programu wypisujacego plik i od skladni kazdej z szesciu powlok

**TerminalProcessSuspendRequest** — Dzis rdzen umie proces wylacznie zakonczyc, wiec dlugie zadanie da sie tylko ubic

**TerminalProcessSuspendResponse** — Dzis rdzen umie proces wylacznie zakonczyc, wiec dlugie zadanie da sie tylko ubic

**TerminalHostSaveRequest** — Bez tej komendy ksiazka zyje jedno posiedzenie przegladarki i ginie przy odswiezeniu strony

**TerminalHostSaveResponse** — Bez tej komendy ksiazka zyje jedno posiedzenie przegladarki i ginie przy odswiezeniu strony

**TerminalHostRemoveRequest** — Karty juz otwarte do tego hosta biegna dalej

**TerminalHostRemoveResponse** — Karty juz otwarte do tego hosta biegna dalej

**TerminalScriptSaveRequest** — Bez tej komendy biblioteka zyje jedno posiedzenie, a jedyna droga jej zachowania jest wywoz do pliku

**TerminalScriptSaveResponse** — Bez tej komendy biblioteka zyje jedno posiedzenie, a jedyna droga jej zachowania jest wywoz do pliku

**TerminalScriptLintRequest** — Analize prowadza programy spoza instalki Danaco Console, wiec odpowiedz mowi wprost, czy narzedzie bylo dostepne — pusty wykaz uwag przy braku narzedzia znaczylby falszywie tresc bez zastrzezen

**TerminalScriptLintResponse** — Analize prowadza programy spoza instalki Danaco Console, wiec odpowiedz mowi wprost, czy narzedzie bylo dostepne — pusty wykaz uwag przy braku narzedzia znaczylby falszywie tresc bez zastrzezen

**TerminalTunnelOpenRequest** — Dzis tunel da sie zalozyc wylacznie poleceniem wydanym w karcie, a wtedy jego stan i przepustowosc sa niewidoczne

**TerminalTunnelOpenResponse** — Dzis tunel da sie zalozyc wylacznie poleceniem wydanym w karcie, a wtedy jego stan i przepustowosc sa niewidoczne

**TerminalKeyGenerateRequest** — Haslo klucza wchodzi odwolaniem do sejfu, nigdy trescia — zgodnie z zasada zapisana w kontrakcie przy zmiennej srodowiska

**TerminalKeyGenerateResponse** — Haslo klucza wchodzi odwolaniem do sejfu, nigdy trescia — zgodnie z zasada zapisana w kontrakcie przy zmiennej srodowiska

**TerminalKeyImportRequest** — Klucz wskazuje sie sciezka, a nie trescia: material kluczowy nie ma powodu przechodzic przez lacze

**TerminalKeyImportResponse** — Klucz wskazuje sie sciezka, a nie trescia: material kluczowy nie ma powodu przechodzic przez lacze

**TerminalKeyRemoveRequest** — Wpisy ksiazki hostow wskazujace ten klucz traca wskazanie i wracaja do klucza domyslnego konfiguracji maszyny

**TerminalKeyRemoveResponse** — Wpisy ksiazki hostow wskazujace ten klucz traca wskazanie i wracaja do klucza domyslnego konfiguracji maszyny

**TerminalWatchStartRequest** — Kontrakt daje dzis wyzwalacz plikowy automatyce, a nie karcie powloki

**TerminalWatchStartResponse** — Kontrakt daje dzis wyzwalacz plikowy automatyce, a nie karcie powloki

**TerminalWatchStopRequest** — Polecenie juz uruchomione biegnie dalej

**TerminalWatchStopResponse** — Polecenie juz uruchomione biegnie dalej

**WorkspaceAgentUnassignRequest** — Nie usuwa eksperta z biblioteki modulu Agents — znosi wylacznie jego przypisanie do tego projektu

**WorkspaceAgentUnassignResponse** — Nie usuwa eksperta z biblioteki modulu Agents — znosi wylacznie jego przypisanie do tego projektu

**WorkspaceAgentUnassignResponse.Unassigned** — Falsz znaczy, ze rdzen takiego przypisania nie mial — nie jest to blad wywolania

**WorkspaceProjectStatusSetRequest** — Droga do stanu 'paused', ktorego opracowanie wymaga, a ktorego zaden dzisiejszy uchwyt nie potrafi zapisac

**WorkspaceProjectStatusSetResponse** — Droga do stanu 'paused', ktorego opracowanie wymaga, a ktorego zaden dzisiejszy uchwyt nie potrafi zapisac

**WorkspaceInstructionsVersionRestoreRequest** — Przywrocenie zaklada wersje nowa o tresci wersji wskazanej — historia nie jest przepisywana

**WorkspaceInstructionsVersionRestoreResponse** — Przywrocenie zaklada wersje nowa o tresci wersji wskazanej — historia nie jest przepisywana

**WorkspaceTaskUpdateRequest** — Pola pominiete zostaja bez zmiany — wywolanie nie jest podmiana calego zadania

**WorkspaceTaskUpdateResponse** — Pola pominiete zostaja bez zmiany — wywolanie nie jest podmiana calego zadania

**WorkspaceTaskMoveRequest** — Osobno od workspace.task.update, bo przeciagniecie karty zmienia stan i porzadek naraz, a jest czynnoscia jednym ruchem myszy

**WorkspaceTaskMoveResponse** — Osobno od workspace.task.update, bo przeciagniecie karty zmienia stan i porzadek naraz, a jest czynnoscia jednym ruchem myszy

**WorkspaceTaskMoveResponse.WipExceeded** — Ostrzezenie, nie odmowa — platforma nie stawia blokad

**WorkspaceTaskDependencySetRequest** — Rdzen odmawia zalozenia zaleznosci domykajacej cykl, bo cykl nie da sie ulozyc w czasie

**WorkspaceTaskDependencySetResponse** — Rdzen odmawia zalozenia zaleznosci domykajacej cykl, bo cykl nie da sie ulozyc w czasie

**WorkspaceScheduleGetRequest** — Zadanie bez obu granic czasu do harmonogramu nie wchodzi

**WorkspaceScheduleGetResponse** — Zadanie bez obu granic czasu do harmonogramu nie wchodzi

**WorkspaceScheduleGetResponse.UnscheduledTaskIds** — Pole niezbedne: bez niego zadania te znikalyby z widoku bez sladu

**WorkspaceCalendarImportResponse.Skipped** — Milczace pomijanie zostawiloby Operatora z niepelnym kalendarzem bez sladu

**WorkspaceNoteSaveRequest** — Puste 'noteId' zaklada notatke nowa; podane zmienia istniejaca. Rdzen przy zapisie przelicza odnosniki tresci

**WorkspaceNoteSaveResponse** — Puste 'noteId' zaklada notatke nowa; podane zmienia istniejaca. Rdzen przy zapisie przelicza odnosniki tresci

**WorkspaceKnowledgeGraphGetRequest** — Nazwa rodziny 'workspace.knowledge' nie miesza sie z rodzina 'knowledge.*': tamta prowadzi wskaznik ZNACZENIA, ta rysuje siec ODNOSNIKOW miedzy bytami projektu

**WorkspaceKnowledgeGraphGetResponse** — Nazwa rodziny 'workspace.knowledge' nie miesza sie z rodzina 'knowledge.*': tamta prowadzi wskaznik ZNACZENIA, ta rysuje siec ODNOSNIKOW miedzy bytami projektu

**WorkspaceCanvasGetResponse.Canvas** — Pole pominiete znaczy, ze projekt nie ma jeszcze zadnej tablicy — nie jest to odmowa

**WorkspaceLibraryTextExtractRequest** — Obejmuje rozpoznanie tekstu z obrazow i skanow

**WorkspaceLibraryTextExtractResponse** — Obejmuje rozpoznanie tekstu z obrazow i skanow

**WorkspaceLibraryDuplicateListRequest** — Sam wykaz niczego nie scala

**WorkspaceLibraryDuplicateListResponse** — Sam wykaz niczego nie scala

**WorkspaceLibraryDuplicateMergeRequest** — Etykiety, kolekcje i wersje plikow scalanych przechodza na plik zachowany, zeby scalenie nie gubilo dorobku

**WorkspaceLibraryDuplicateMergeResponse** — Etykiety, kolekcje i wersje plikow scalanych przechodza na plik zachowany, zeby scalenie nie gubilo dorobku

**WorkspaceSearchProjectRequest** — Nie zastepuje library.file.search: tamta komenda przeszukuje biblioteke centralna, ta jeden projekt i wiecej niz pliki. Wyszukiwanie po ZNACZENIU prowadzi rodzina knowledge.

**WorkspaceSearchProjectResponse** — Nie zastepuje library.file.search: tamta komenda przeszukuje biblioteke centralna, ta jeden projekt i wiecej niz pliki. Wyszukiwanie po ZNACZENIU prowadzi rodzina knowledge.

**LibraryMetadataGetRequest** — Metadane osadzone czyta sie z bajtow, wiec ich odczyt jest kosztowny i wchodzi wylacznie na wyrazne zadanie

**LibraryMetadataGetResponse** — Metadane osadzone czyta sie z bajtow, wiec ich odczyt jest kosztowny i wchodzi wylacznie na wyrazne zadanie

**LibraryMetadataSetRequest** — Dzis formularz opisu nie ma dokad pojsc: library.tag.set zmienia wylacznie etykiety i kolekcje. Zapis rozglasza library.file.changed

**LibraryMetadataSetResponse** — Dzis formularz opisu nie ma dokad pojsc: library.tag.set zmienia wylacznie etykiety i kolekcje. Zapis rozglasza library.file.changed

**LibrarySchemaSetRequest** — Pole o kodzie juz istniejacym jest zmieniane, nie dublowane

**LibrarySchemaSetResponse** — Pole o kodzie juz istniejacym jest zmieniane, nie dublowane

**LibraryTagListRequest** — Bez tej komendy klient sklada slownik z etykiet plikow odczytanej strony wykazu, wiec pokazuje probke zamiast slownika

**LibraryTagListResponse** — Bez tej komendy klient sklada slownik z etykiet plikow odczytanej strony wykazu, wiec pokazuje probke zamiast slownika

**LibraryTagUpdateRequest** — Zmiana nazwy przechodzi na wszystkie zasoby noszace etykiete — inaczej powstalaby druga etykieta o tym samym znaczeniu

**LibraryTagUpdateResponse** — Zmiana nazwy przechodzi na wszystkie zasoby noszace etykiete — inaczej powstalaby druga etykieta o tym samym znaczeniu

**LibraryTagMergeRequest** — Zasoby noszace etykiete zrodlowa dostaja docelowa, a zrodlowa znika ze slownika

**LibraryTagMergeResponse** — Zasoby noszace etykiete zrodlowa dostaja docelowa, a zrodlowa znika ze slownika

**LibraryTagRemoveRequest** — Usuniecie etykiety uzywanej zada potwierdzenia, bo zdejmuje ja z zasobow, ktorych zadanie nie wymienia

**LibraryTagRemoveResponse** — Usuniecie etykiety uzywanej zada potwierdzenia, bo zdejmuje ja z zasobow, ktorych zadanie nie wymienia

**LibraryCollectionListRequest** — Bez tej komendy klient zna wylacznie identyfikatory kolekcji wyczytane z plikow, a nazwe kolekcji tylko w chwili jej zalozenia

**LibraryCollectionListResponse** — Bez tej komendy klient zna wylacznie identyfikatory kolekcji wyczytane z plikow, a nazwe kolekcji tylko w chwili jej zalozenia

**LibraryRuleSetRequest** — Zapis uruchamia przeliczenie

**LibraryRuleSetResponse** — Zapis uruchamia przeliczenie

**LibraryRuleRemoveRequest** — Zasoby przypisane przez regule kolekcji inteligentnej zostaja w kolekcji, ale przestaja byc przeliczane — inaczej usuniecie reguly oproznialoby kolekcje po cichu

**LibraryRuleRemoveResponse** — Zasoby przypisane przez regule kolekcji inteligentnej zostaja w kolekcji, ale przestaja byc przeliczane — inaczej usuniecie reguly oproznialoby kolekcje po cichu

**LibraryDuplicateScanRequest** — Rozpoznanie przyblizone zada porownania calego zbioru, wiec jest czynnoscia rdzenia, nie okna

**LibraryDuplicateScanResponse** — Rozpoznanie przyblizone zada porownania calego zbioru, wiec jest czynnoscia rdzenia, nie okna

**LibraryDuplicateMergeRequest** — Zasoby wchloniete trafiaja do archiwum, nie znikaja

**LibraryDuplicateMergeResponse** — Zasoby wchloniete trafiaja do archiwum, nie znikaja

**LibraryFixityCheckRequest** — Klient nie ma czym tego zrobic, bo nie siega po bajty zasobu

**LibraryFixityCheckResponse** — Klient nie ma czym tego zrobic, bo nie siega po bajty zasobu

**LibraryNameNormalizeRequest** — Przebieg probny pokazuje wynik bez zapisu

**LibraryNameNormalizeResponse** — Przebieg probny pokazuje wynik bez zapisu

**LibraryAuditListRequest** — Dziennik jest przyrostowy — wpisow nie da sie zmienic ani usunac ta rodzina komend

**LibraryAuditListResponse** — Dziennik jest przyrostowy — wpisow nie da sie zmienic ani usunac ta rodzina komend

**LibraryRetentionSetRequest** — Twarde usuniecie nie zachodzi samoczynnie: polityka najwyzej zglasza zasob do usuniecia

**LibraryRetentionSetResponse** — Twarde usuniecie nie zachodzi samoczynnie: polityka najwyzej zglasza zasob do usuniecia

**LibraryPreservationRunRequest.Profile** — PDF/A-2b; brak bierze odmiane z konfiguracji modulu

**LibraryFileMoveRequest** — Przeniesienie nie rusza tresci ani historii wersji

**LibraryFileMoveResponse** — Przeniesienie nie rusza tresci ani historii wersji

**LibraryFileArchiveRequest** — Czynnosc jest odwracalna komenda library.file.restore — to jest kosz repozytorium, nie usuniecie

**LibraryFileArchiveResponse** — Czynnosc jest odwracalna komenda library.file.restore — to jest kosz repozytorium, nie usuniecie

**LibraryFileDeleteRequest** — Jedyna droga utraty danych repozytorium; bez potwierdzenia jest odmowa

**LibraryFileDeleteResponse** — Jedyna droga utraty danych repozytorium; bez potwierdzenia jest odmowa

**LibraryShareCreateRequest** — Token jest jawny zgodnie z zasada jawnosci kluczy platformy

**LibraryShareCreateResponse** — Token jest jawny zgodnie z zasada jawnosci kluczy platformy

**LibraryShareRevokeRequest** — Odnosnik przestaje dzialac, wpis zostaje w dzienniku audytu

**LibraryShareRevokeResponse** — Odnosnik przestaje dzialac, wpis zostaje w dzienniku audytu

**LibraryClassifyRunRequest** — Wynik wraca sugestiami do przyjecia, nie zapisem — domyslnym zachowaniem modulu jest sugestia z akceptacja Operatora

**LibraryClassifyRunResponse** — Wynik wraca sugestiami do przyjecia, nie zapisem — domyslnym zachowaniem modulu jest sugestia z akceptacja Operatora

**LibrarySuggestionApplyRequest** — Przyjecie wykonuje czynnosc, ktora sugestia opisuje; odrzucenie zdejmuje ja z wykazu

**LibrarySuggestionApplyResponse** — Przyjecie wykonuje czynnosc, ktora sugestia opisuje; odrzucenie zdejmuje ja z wykazu

**LibraryDiffCompareRequest** — Dokumenty binarne porownuje po tekscie z nich wydobytym i mowi o tym w odpowiedzi

**LibraryDiffCompareResponse** — Dokumenty binarne porownuje po tekscie z nich wydobytym i mowi o tym w odpowiedzi

**DesignVectorPathSetRequest** — Wezly przychodza w calosci — zmiana jednego wezla idzie ta sama droga co narysowanie sciezki, zeby nie bylo dwoch prawd o jej ksztalcie

**DesignVectorPathSetResponse** — Wezly przychodza w calosci — zmiana jednego wezla idzie ta sama droga co narysowanie sciezki, zeby nie bylo dwoch prawd o jej ksztalcie

**DesignVectorShapeAddRequest** — Ksztalt powstaje od razu jako wezly, wiec da sie go dalej edytowac pioram — nie jest osobnym bytem, ktory potem trzeba zamieniac

**DesignVectorShapeAddResponse** — Ksztalt powstaje od razu jako wezly, wiec da sie go dalej edytowac pioram — nie jest osobnym bytem, ktory potem trzeba zamieniac

**DesignVectorBooleanRequest** — Sciezki zrodlowe znikaja albo zostaja wedle wskazania — bo suma dwoch ksztaltow bywa krokiem posrednim, a bywa wynikiem koncowym

**DesignVectorBooleanResponse** — Sciezki zrodlowe znikaja albo zostaja wedle wskazania — bo suma dwoch ksztaltow bywa krokiem posrednim, a bywa wynikiem koncowym

**DesignVectorTextPathRequest** — Zamiana w kontury jest nieodwracalna dla wyniku, dlatego tekst zrodlowy zostaje

**DesignVectorTextPathResponse** — Zamiana w kontury jest nieodwracalna dla wyniku, dlatego tekst zrodlowy zostaje

**DesignVectorOptimizeRequest** — Odpowiedz niesie ubytek zmierzony, zeby Operator widzial, ile naprawde ubylo, zamiast czytac obietnice

**DesignVectorOptimizeResponse** — Odpowiedz niesie ubytek zmierzony, zeby Operator widzial, ile naprawde ubylo, zamiast czytac obietnice

**DesignVectorSymbolSetRequest** — Zmiana definicji propaguje do wszystkich instancji — po to symbol jest

**DesignVectorSymbolSetResponse** — Zmiana definicji propaguje do wszystkich instancji — po to symbol jest

**DesignFrameSetRequest** — Ramka jest ekranem makiety: to ona, a nie kanwa, wyznacza obszar wydania

**DesignFrameSetResponse** — Ramka jest ekranem makiety: to ona, a nie kanwa, wyznacza obszar wydania

**DesignFrameRemoveRequest** — Warstwy ramki zostaja na kanwie — usuniecie ramki nie jest usunieciem pracy, ktora w niej lezala

**DesignFrameRemoveResponse** — Warstwy ramki zostaja na kanwie — usuniecie ramki nie jest usunieciem pracy, ktora w niej lezala

**DesignLayoutAutoRequest** — Rdzen oddaje warstwy po przeliczeniu, wiec klient nie liczy ukladu drugi raz i nie ma jak sie rozjechac

**DesignLayoutAutoResponse** — Rdzen oddaje warstwy po przeliczeniu, wiec klient nie liczy ukladu drugi raz i nie ma jak sie rozjechac

**DesignFrameResizeApplyRequest** — Osobno od design.frame.set, bo tam zmiana rozmiaru jest zapisem nastawy, a tu przeliczeniem ukladu wedle wiezi

**DesignFrameResizeApplyResponse** — Osobno od design.frame.set, bo tam zmiana rozmiaru jest zapisem nastawy, a tu przeliczeniem ukladu wedle wiezi

**DesignMockupGenerateRequest** — Wynik jest makieta do poprawienia, nie projektem koncowym, i tak sie o nim mowi

**DesignMockupGenerateResponse** — Wynik jest makieta do poprawienia, nie projektem koncowym, i tak sie o nim mowi

**DesignMockupImportRequest** — Odpowiedz mowi, czego rdzen nie rozpoznal — bilans zamiast ciszy

**DesignMockupImportResponse** — Odpowiedz mowi, czego rdzen nie rozpoznal — bilans zamiast ciszy

**DesignColorVisionSimulateRequest** — Wynik jest nowym zasobem — zrodlo zostaje nietkniete

**DesignColorVisionSimulateResponse** — Wynik jest nowym zasobem — zrodlo zostaje nietkniete

**DesignColorAccessibilityAuditRequest** — Czynnosc CZYTA i niczego nie zmienia

**DesignColorAccessibilityAuditResponse** — Czynnosc CZYTA i niczego nie zmienia

**DesignStockSearchRequest** — Brak skonfigurowanego dostawcy jest ODPOWIEDZIA nazywajaca brak, nie cisza

**DesignStockSearchResponse** — Brak skonfigurowanego dostawcy jest ODPOWIEDZIA nazywajaca brak, nie cisza

**DesignPrintPreflightRequest** — Czynnosc CZYTA i niczego nie zmienia; zastrzezenie z waga blad zatrzyma druk, wiec Operator ma je zobaczyc TUTAJ, a nie w drukarni

**DesignPrintPreflightResponse** — Czynnosc CZYTA i niczego nie zmienia; zastrzezenie z waga blad zatrzyma druk, wiec Operator ma je zobaczyc TUTAJ, a nie w drukarni

**DesignLargeformatTileRequest** — Odpowiedz niesie uklad kafli, wiec Operator wie, ktory kafel gdzie idzie, zanim cokolwiek wydrukuje

**DesignLargeformatTileResponse** — Odpowiedz niesie uklad kafli, wiec Operator wie, ktory kafel gdzie idzie, zanim cokolwiek wydrukuje

**DesignChartRenderRequest** — Dane przychodza seriami, a nie obrazem — wykres da sie wiec przerysowac po zmianie liczb, zamiast rysowac go od nowa

**DesignChartRenderResponse** — Dane przychodza seriami, a nie obrazem — wykres da sie wiec przerysowac po zmianie liczb, zamiast rysowac go od nowa

**DesignPhotoCropRequest** — Wynik jest wariantem zrodla — oryginal zostaje nietkniety

**DesignPhotoCropResponse** — Wynik jest wariantem zrodla — oryginal zostaje nietkniety

**DesignPhotoUpscaleRequest** — Gdy kanal modelu obrazowego stoi, liczy kanal; gdy nie stoi, liczy rachunek wkompilowany — odpowiedz mowi to polem computedBy

**DesignPhotoUpscaleResponse** — Gdy kanal modelu obrazowego stoi, liczy kanal; gdy nie stoi, liczy rachunek wkompilowany — odpowiedz mowi to polem computedBy

**DesignPhotoInpaintRequest** — Gdy kanal modelu obrazowego stoi, liczy kanal; gdy nie stoi, liczy rachunek wkompilowany — odpowiedz mowi to polem computedBy

**DesignPhotoInpaintResponse** — Gdy kanal modelu obrazowego stoi, liczy kanal; gdy nie stoi, liczy rachunek wkompilowany — odpowiedz mowi to polem computedBy

**DesignPhotoExpandRequest** — Gdy kanal modelu obrazowego stoi, liczy kanal; gdy nie stoi, liczy rachunek wkompilowany — odpowiedz mowi to polem computedBy

**DesignPhotoExpandResponse** — Gdy kanal modelu obrazowego stoi, liczy kanal; gdy nie stoi, liczy rachunek wkompilowany — odpowiedz mowi to polem computedBy

**DesignPhotoBackgroundRemoveRequest** — Gdy kanal modelu obrazowego stoi, liczy kanal; gdy nie stoi, liczy rachunek wkompilowany — odpowiedz mowi to polem computedBy

**DesignPhotoBackgroundRemoveResponse** — Gdy kanal modelu obrazowego stoi, liczy kanal; gdy nie stoi, liczy rachunek wkompilowany — odpowiedz mowi to polem computedBy

**AodMuteGetRequest** — Wyciszenie przeterminowane nie wchodzi do wykazu — Operator nie ma go odklikiwac

**AodMuteGetResponse** — Wyciszenie przeterminowane nie wchodzi do wykazu — Operator nie ma go odklikiwac

**AodMuteSetRequest** — Zniesienie idzie ta sama komenda z polem muted rownym falszowi — odwracalnosc jednym ruchem jest wymogiem, nie wygoda. Odpowiedz oddaje wykaz wyciszen PO zmianie, zeby powloka nie musiala pytac drugi raz. Zmiana rozglasza sie zdarzeniem aod.mute.changed na pozostale powloki Operatora

**AodMuteSetRequest.MuteId** — Podane przy muted rownym falszowi znosi dokladnie to wyciszenie; puste znosi wyciszenie zlozone z rodzaju i zakresu tego zadania

**AodMuteSetRequest.EndsAt** — Wyciszenie czasowe bez niej jest zadaniem niezgodnym z kontraktem — cisza bez konca nie mowi Operatorowi, do kiedy trwa

**AodMuteSetResponse** — Zniesienie idzie ta sama komenda z polem muted rownym falszowi — odwracalnosc jednym ruchem jest wymogiem, nie wygoda. Odpowiedz oddaje wykaz wyciszen PO zmianie, zeby powloka nie musiala pytac drugi raz. Zmiana rozglasza sie zdarzeniem aod.mute.changed na pozostale powloki Operatora

**AodMuteSetResponse.Changed** — Falsz znaczy wyciszenie powtorzone albo zniesienie wyciszenia, ktorego nie bylo

**AodSignalReportRequest** — Nosnik dla trzech klas rozdz. 3.2, ktorych rdzen nie widzi wlasna telemetria: wyniku kontroli jakosci, harmonogramu przebiegow automatyk i powtarzalnosci czynnosci Operatora. Odpowiedz mowi wprost, czy sygnal wpadl w wyciszenie i ktore — sygnal wyciszony odklada sie nadal, bo wyciszenie wstrzymuje UJAWNIENIE, a nie zapis

**AodSignalReportResponse** — Nosnik dla trzech klas rozdz. 3.2, ktorych rdzen nie widzi wlasna telemetria: wyniku kontroli jakosci, harmonogramu przebiegow automatyk i powtarzalnosci czynnosci Operatora. Odpowiedz mowi wprost, czy sygnal wpadl w wyciszenie i ktore — sygnal wyciszony odklada sie nadal, bo wyciszenie wstrzymuje UJAWNIENIE, a nie zapis

**AodSignalListRequest** — Odpowiedz nazywa wyciszenia, ktore sygnaly wstrzymaly, i liczbe wstrzymanych — cisza, po ktorej Operator nie wie, ze cos milczy, jest gorsza od braku wyciszenia

**AodSignalListResponse** — Odpowiedz nazywa wyciszenia, ktore sygnaly wstrzymaly, i liczbe wstrzymanych — cisza, po ktorej Operator nie wie, ze cos milczy, jest gorsza od braku wyciszenia

**MemoryDisableListRequest** — Odczyt do pary z memory.disable.set; okno konfiguracji buduje z niego zakres ustawien „Pamiec"

**MemoryDisableListResponse** — Odczyt do pary z memory.disable.set; okno konfiguracji buduje z niego zakres ustawien „Pamiec"

**MemoryDisableSetRequest** — Zniesienie idzie ta sama komenda z polem disabled rownym falszowi: odwracalnosc jednym ruchem. WYLACZENIE NIE KASUJE TRESCI — tym rozni sie od memory.delete; wpis wylaczony zostaje na miejscu i wraca w calosci po zniesieniu. Odpowiedz oddaje wykaz wylaczen PO zmianie

**MemoryDisableSetResponse** — Zniesienie idzie ta sama komenda z polem disabled rownym falszowi: odwracalnosc jednym ruchem. WYLACZENIE NIE KASUJE TRESCI — tym rozni sie od memory.delete; wpis wylaczony zostaje na miejscu i wraca w calosci po zniesieniu. Odpowiedz oddaje wykaz wylaczen PO zmianie

**MemoryDisableSetResponse.Changed** — Falsz znaczy wylaczenie powtorzone albo zniesienie wylaczenia, ktorego nie bylo

**NotificationListRequest** — Filtry zawezaja wykaz do klas, wag, stanow albo zrodla; pominiete zwracaja rejestr niezamkniety. Licznik zdarzen nowych idzie zawsze, niezaleznie od filtru — plakietka paska kontekstu liczy caly rejestr, nie widok

**NotificationListResponse** — Filtry zawezaja wykaz do klas, wag, stanow albo zrodla; pominiete zwracaja rejestr niezamkniety. Licznik zdarzen nowych idzie zawsze, niezaleznie od filtru — plakietka paska kontekstu liczy caly rejestr, nie widok

**NotificationAcknowledgeRequest** — Wykaz pusty oznacza wszystkie zdarzenia nowe, czyli dzialanie zbiorcze centrum

**NotificationAcknowledgeResponse** — Wykaz pusty oznacza wszystkie zdarzenia nowe, czyli dzialanie zbiorcze centrum

**NotificationResolveRequest** — Dzialanie, ktore je zamyka, wykonuje rodzina wlasciwa — ta komenda zapisuje wylacznie skutek w rejestrze centrum

**NotificationResolveResponse** — Dzialanie, ktore je zamyka, wykonuje rodzina wlasciwa — ta komenda zapisuje wylacznie skutek w rejestrze centrum

**NotificationSnoozeRequest** — Po niej zdarzenie wraca do stanu nowe i znow liczy sie do plakietki — odlozenie nie jest zamknieciem

**NotificationSnoozeResponse** — Po niej zdarzenie wraca do stanu nowe i znow liczy sie do plakietki — odlozenie nie jest zamknieciem

**StreamChunkEvent.MessageId** — Znaczenie pola jest STALE od pierwszego fragmentu do domykajacego — nadawca nie ma prawa zmienic go w polowie strumienia. W debacie jest to kod WYPOWIEDZI, nie kod uczestnika; mowce klient bierze z RoundtableStatement.participantId

**AuthChangedEvent** — Sekcja Uwierzytelnianie Okna Ustawien odswieza sie bez odpytywania

**DeviceChangedEvent** — Rozglaszane do WSZYSTKICH polaczonych urzadzen, zeby uniewaznienie jednego bylo natychmiast widoczne na pozostalych ekranach — tej samej drogi uzywa odzyskanie konta, ktore uniewaznia tokeny wydane przed zmiana hasla

**AdvisorConsultedEvent** — Zdarzenie istnieje dla JAWNOSCI: rada, ktorej nie widac w strumieniu, jest przeslanka, ktorej w aktach nie ma — a agent na niej dziala. Niesie okno, doradce, podstawe doboru i skrot rady; pelna tresc wraca wynikiem komendy advisor.consult i zapisem dziennika konsultacji

**AppsWorkspaceChangedEvent** — Bez tego zdarzenia apps.workspace.update byl jedyna komenda zapisu modulu, ktora nic nie rozglaszala — drugie okno nie dowiadywalo sie o zmianie

**SessionToolAttachedEvent** — Skoro model dostaje narzedzie, ktorego nie mial, Operator MUSI to zobaczyc — kazde posuniecie widac na ekranie

**SessionToolDetachedEvent** — Widocznosc obowiazuje w OBIE strony: zestaw, ktory urosl na oczach Operatora, nie moze skurczyc sie po cichu — bez tego zdarzenia sasiednie okno pokazywaloby narzedzie, ktorego model juz nie ma

**AppsPreviewChangedEvent** — Bez tego zdarzenia podglad na zywo bylby odswiezaniem na zadanie, czyli nie podgladem na zywo

**AppsDeploymentLogEvent** — Osobne od apps.build.changed, ktore niesie STAN wdrozenia, a nie jego dziennik

**SpeechListenPartialEvent** — Tekst jest nietrwaly i bywa poprawiany kolejnym zdarzeniem — trwaly zapis daje dopiero speech.transcribe

**AutomationExecutionLoggedEvent** — Nazwa rozna od komendy automation.execution.log, bo kontrakt nie dopuszcza tej samej nazwy w obu wykazach

**DesignBoardChangedEvent** — Dzis kompozycja jezdzi w obie strony komendami, ale drugie okno i drugie polaczenie nie dowiaduja sie o zmianie niczym — zdarzenia obszaru jest dokladnie jedno i dotyczy zasobu

**DesignBoardPresenceEvent** — Zdarzenie ULOTNE — nie ma odpowiednika w bazie i nie jest odtwarzane po ponownym polaczeniu

**AlertTriggeredEvent** — Zdarzenie jest droga alertu do Always On Display i do powiadomienia w aplikacji; rejestr odpytywany komenda alert.trigger.list sam by tam nie dotarl, bo alert ma znaczenie wtedy, gdy Operator nie patrzy

**AodMuteChangedEvent** — Rozgloszenie sciga cisza pozostale powloki Operatora: bez niego Operator wyciszalby w jednej, a sugestie wchodzilyby w drugiej

**MemoryDisableChangedEvent** — Rozgloszenie idzie na pozostale powloki Operatora, zeby wykaz pamieci w kazdej z nich pokazywal ten sam stan wlaczenia

**NotificationRaisedEvent** — Rozglaszane do wszystkich polaczen Operatora, zeby plakietka i kolumna centrum nie musialy odpytywac rdzenia w tle; to samo zdarzenie jest podstawa Toastu w chwili wystapienia

**NotificationChangedEvent** — Rozglaszane do wszystkich polaczen, zeby obsluga na jednym urzadzeniu byla natychmiast widoczna na pozostalych

**UnknownCommandPayload** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**ZdarzenieNieznanej** — Fail-open: polaczenie nie jest zrywane, a obszar spoza kontraktu dostaje zdarzenie zapasowe.

**NarzedziaModelu** — Kanal modelu dostaje sterowanie platforma jako narzedzia, a nie jako osobny parser intencji; kazda deklaracja odpowiada komendzie kontraktu.

**KnownModuleIds** — Lista informacyjna, nie brama.

**KnownChannelKinds** — Lista informacyjna; nowy kanal = nowy wiersz rejestru, nie zmiana kodu.

**KnownSettingKeys** — Lista informacyjna, nie brama — katalog ustawien jest sterowany danymi, wiec nowa pozycja to nowy wiersz, nie zmiana kodu.

**KnownSessionConfigKeys** — Klucz powstaje mechanicznie: "sesja.konfiguracja." + wartosc SessionConfigArea. Konfiguracja sesji nie ma wlasnego rezolwera — idzie tym samym rezolwerem osmiu poziomow i trzech osi, ktory obsluguje config.get i config.set. Lista informacyjna, nie brama.

### budowa/server/internal/dane/studio_kontrola_pracy.go

Uzasadnienia i zastrzeżenia projektowe przeniesione z komentarzy warstwy danych blokad, dziennika czynności i kopii zapasowych modułu Studio.

**Plik** — Trzy obszary (blokady fragmentów, dziennik czynności wraz z zależnościami, kopie zapasowe i nastawy pracy) leżą w jednym pliku, bo wiąże je jedno pytanie zadawane w jednym miejscu: czy tę zmianę wolno wnieść, a jeśli tak, to czym ją potem cofnąć. Blokada odpowiada na pierwszą połowę, dziennik na drugą, a kopia zapasowa jest siatką pod obiema — zakłada się ją przed czynnością nieodwracalną i po nieudanym zapisie. Rozdzielenie ich na trzy pliki rozdzieliłoby zapytania, które i tak padają razem. Znakowanie fragmentów, zajęcia wykonawców i spięcia leżą osobno, w warstwie danych znakowania wykonawcy: tamte opisują, co ktoś o dokumencie powiedział i kto nad nim pracuje, a nie czego nie wolno tknąć. Wygasanie zajęcia i wygasanie kopii zapasowej liczy się w SQL (`strftime('now')`), bo obie wielkości muszą być tym samym zegarem co kolumny `utworzono` i `wygasa` w tych wierszach — zegar rdzenia i zegar bazy rozjadą się przy pierwszej różnicy strefy, a wtedy kopia wygasłaby wcześniej albo później, niż mówi nastawa Operatora.

**BlokadaFragmentuStudia** — Pole `Zasieg` rozstrzyga, kogo blokada dotyczy: wartość `model` (postać domyślna) wiąże wyłącznie wykonawców, wartość `everyone` wiąże także Operatora. Operator zmienia fragment zablokowany bez przeszkód, dopóki nie zażąda zasięgu `everyone` jawnie.

**CzynnoscDokumentuStudia** — Pola `StanPrzed` i `StanPo` niosą wycinek objęty czynnością, nie migawkę całego dokumentu — inaczej cofnięcie czynności ze środka dziennika zabrałoby ze sobą wszystko, co po niej weszło, czyli byłoby przywróceniem wersji zamiast cofnięciem jednej zmiany.

**KopiaZapasowaStudia** — Pola `UdaloSie` i `PowodNiepowodzenia` istnieją, bo wskaźnik zapisano pokazany przy zapisie nieudanym jest najgorszym możliwym błędem tego modułu: Operator zamknie okno i straci pracę. Kopia nieudana zostaje wierszem — musi być widoczna i nazwana, a nie zniknąć razem z niepowodzeniem.

**blokadyStudiaPrzesun** — Przesunięcie zakresu po wpisie następuje przed blokadą leżącą za punktem edycji. Bez tego blokada zaczęłaby po pierwszej edycji chronić nie ten fragment, co miała, i to bez śladu. Kolumna `przesuniecia` rośnie razem z zakresem: jest miarą zaufania do zakresu, widoczną dla Operatora w wykazie.

**kopieStudiaPrzemiec** — Zasada wygasania kopii jest jawnym, odwracalnym ustawieniem Operatora, więc granice podaje wołający, a nie stała rdzenia. Kopia nieudana nie wygasa razem z udanymi: jest jedynym śladem, że praca nie doszła na dysk, i ma zostać, dopóki Operator jej nie zobaczy.

**PrzesunBlokadyFragmentow** — Blokada leżąca przed punktem edycji zostaje nietknięta; leżąca za nim jedzie o różnicę długości. Blokada, w której środek trafiła edycja, też zostaje nietknięta — skoro edycja weszła w blokadę, wolno jej było tam wejść (Operator albo zasięg `model` przy czynności Operatora), a wtedy zakres blokady ma zostać taki, jaki Operator ustawił.

**ZapiszCzynnoscDokumentu** — Kolejność nadaje baza, nie wołający: dwie czynności zapisane w tej samej chwili muszą dostać różne numery, a numer nadany przez rdzeń z odczytu stanu bieżącego byłby wyścigiem. Więz UNIQUE(dokument, kolejnosc) zamienia ten wyścig w błąd zapisu, zamiast w dwa wpisy o tym samym miejscu w porządku.

**CzynnosciDokumentu** — Zależności doczytywane są jednym zapytaniem dla całego dokumentu, nie zapytaniem na wpis: dziennik długiego dokumentu ma setki wpisów, a pytanie na każdy wpis zamieniłoby odczyt wykazu w setki zapytań.

**NastawaPracy** — Wiersz nastaw zakłada się przy odczycie, nie przy zapisie, bo odczyt nastaw autozapisu ma oddać nastawy obowiązujące także wtedy, gdy Operator nigdy ich nie ruszał — a wtedy obowiązują wartości domyślne kolumn, i to jest odpowiedź, nie brak danych.

**ZapiszSkutekAutozapisu** — Skutek zapisu odkłada się osobnym poleceniem, bo pisze go inna czynność niż nastawy: nastawy stawia Operator, skutek zapisuje sam mechanizm autozapisu. Jedno wspólne polecenie kazałoby autozapisowi przepisywać nastawy, których nie zmieniał.

### budowa/server/internal/dane/agent_zakres.go

Uzasadnienia i zastrzeżenia projektowe przeniesione z komentarzy zakresu działania eksperta widzianego od strony eksperta.

**Plik** — Repozytorium osobne, a nie kolejne czynności `RepozytoriumAgentow`: tamten kontrakt opisuje bibliotekę ekspertów — założenie, wykaz, zmianę, usunięcie — i pracuje w nim równolegle więcej rąk. Zakres działania jest inną odpowiedzialnością: odpowiada nie na pytanie „jaki jest ten ekspert", lecz „co temu ekspertowi wolno zrobić w systemie". Zapis, którego nikt nie czyta przy wykonaniu, jest gorszy niż jego brak — wiersze zakładane tutaj czyta straż zakresu eksperta (`core/straz_eksperta.go`) i na ich podstawie odmawia: nałożenia eksperta na okno modułu, którego nie ma w jego zakresie, oraz powołania podagentów przez eksperta z wyłączonym Subagent Network. Bez tego Operator widziałby ograniczenie, którego nikt nie egzekwuje. Brak wiersza znaczy stan wyjściowy platformy, czyli pełny dostęp: brak modułów znaczy „wszędzie", brak przełącznika izolacji znaczy „nieodcięty".

**PrzypisaniaEkspertow** — Dwa zapytania, nie jedno z UNION: źródła mają inne kolumny i inne znaczenie celu — projekt ma kod własny, stanowisko obsady ma klucz wiersza — a złożenie ich w jedno zapytanie kazałoby czytelnikowi zgadywać, skąd wziął się wiersz.

### budowa/server/internal/dane/studio_znakowanie_wykonawcy.go

Uzasadnienia i zastrzeżenia projektowe przeniesione z komentarzy części drugiej warstwy danych kontroli pracy modułu Studio.

**Plik, wersje w szeregach** — `studio_wersje.go` zna wersję sprzed dobudowy: treść, etykietę, kamień milowy. Kolumny `szereg`, `postac_json` i tożsamości wykonawcy dołożyła migracja 367 i pyta o nie wyłącznie ten odcinek — wykaz historii z przełącznikiem „pokaż także zapisy samoczynne" oraz porównanie postaci dwóch wersji. Dopisanie ich do tamtego pliku byłoby wejściem w plik cudzego odcinka; osobny odczyt tych samych wierszy nie zakłada drugiego pojęcia wersji, bo tabela jest jedna i przywraca się ją tą samą drogą.

**Plik, zajęcie a blokada** — Blokada Operatora jest trwała i skierowana przeciw wykonawcom: zdejmuje ją wyłącznie Operator. Zajęcie fragmentu jest chwilowe, wygasa samo i chroni przed drugim wykonawcą. Dwa różne byty, dwie tabele — zlanie ich dałoby blokadę, która wygasa (czyli żadną), albo zajęcie, którego nikt nie zdejmie po agencie ubitym w pół pracy.

**ZajecieFragmentuStudia.Wygasle** — Liczone zapytaniem SQL, nie w rdzeniu: zegar rdzenia i zegar bazy rozjadą się przy pierwszej różnicy strefy, a wtedy zajęcie trzymałoby fragment dłużej albo krócej, niż mówi jego własna kolumna `wygasa`.

**wersjaZalozycielskaStudia** — Szereg autozapisu odpada z rachunku z zamysłem: gdyby pierwszy zapis samoczynny wypadł przed pierwszym zapisem Operatora, powrót do stanu pierwotnego wracałby do stanu przypadkowego, a nie do tego, co Operator założył.

**ZapiszZajecieFragmentu** — Odstęp zerowy albo ujemny znaczy zajęcie bez wygasania — i takiego zajęcia rdzeń nie zakłada, bo agent ubity w pół pracy trzymałby fragment na zawsze; granicę podaje wołający.

**ZapiszWersjeSzeregu** — Wskaźnika wersji bieżącej nie przestawia, tak samo jak `ZapiszWersje` z `studio_wersje.go`; decyzja, kiedy nowa wersja staje się bieżącą, należy do wołającego.

### budowa/server/internal/dane/badania_dobudowa.go

Uzasadnienia i zastrzeżenia projektowe przeniesione z komentarzy dobudowy obszaru Research po stronie danych.

**Plik** — Ustalenia i ich otoczenie leżą w `badania_dobudowa_ustalenia.go`, odkrywanie i przestrzeń w `badania_dobudowa_odkrycia.go`, raport i eksport w `badania_dobudowa_raport.go`. Wszystkie cztery pliki niosą metody tego samego `*repozytoriumBadan`, tak jak trzy pliki zastane. Kontrakt dobudowy jest jednym interfejsem wpiętym do `RepozytoriumBadan` przez osadzenie: jedna deklaracja, jedno miejsce do przeczytania, a plik `badania.go` nie rośnie o siedemdziesiąt sygnatur, których nie implementuje.

## budowa/server/internal/dane/urzadzenia_powiadomien.go

Rozbicie tabel `urzadzenie_powiadomien`, `powiadomienie` i `powiadomienie_dostarczenie`
na trzy repozytoria zmusiłoby silnik do składania transakcji z kawałków trzech
właścicieli, a najważniejsza droga tego pliku (Takt) jest transakcją obejmującą
wszystkie trzy tabele naraz.

Konstruktor `NowePowiadomienia` jest eksportowany i bierze uchwyt puli, żeby
pakiet silnika powiadomień mógł złożyć repozytorium z uchwytu, który już dostał
w kompozycji, bez dodawania pola w `dane/zestaw.go`. Pula połączeń jest jedna na
proces; to repozytorium jej nie otwiera i nie zamyka.

`KluczKanalu` w `RejestracjaPowiadomien` dla kanału `polaczenie` niesie ten sam
napis, który transport przekazuje jako `Tozsamosc.IdKlienta` — adres gniazda,
nie odrębną tożsamość urządzenia.

`Wyslij` w `PlanTaktu` jest domknięciem, a nie zwracaną wartością, żeby sprawdzenie
stanu wiersza i sama wysyłka zmieściły się w jednej transakcji zapisu: sprawdzenie
przy pobraniu należnych nie wystarcza, bo między pobraniem a wysłaniem decyzja
może zapaść i urządzenie dostałoby powiadomienie po fakcie. Pusta lista przyjętych
odbiorców znaczy „nie było komu" i jest powodem ponowienia, nie błędem.

Kolejność kroków w `Takt` nie może zostać przestawiona: `zajmijPowiadomienie`
bierze zamek zapisu i w tym samym poleceniu sprawdza, że wiersz nadal jest
`oczekuje` — zero zmienionych wierszy znaczy, że decyzja już zapadła, termin minął
albo doręczenie już było. Dopiero pod tym zamkiem czytany jest wiersz i sprawdzany
termin ważności. Wysyłka idzie wewnątrz transakcji, więc koperta wychodzi przy
trzymanym zamku zapisu — cena możliwa do przyjęcia, bo doręczenie kanałem
`polaczenie` jest zapisem do gniazda w pamięci procesu. Odwrotna kolejność —
zwolnić zamek, potem wysłać — przywraca szczelinę, w której powiadomienie
o zapadłej już decyzji zdąży wyjść. Wołający, który nie odda `Wyslij` albo
`NastepnaProba`, dostaje błąd zamiast cichego pominięcia.

Zapytanie `Nalezne` porządkuje wynik tak, że bez tego kolejka pilna stałaby za
zwykłą, która akurat weszła pierwsza, i priorytet nie znaczyłby nic.

## budowa/server/internal/dane/studio_postac_dokumentu.go

Plik deklaruje kontrakt obszaru postaci — drzewa postaci, arkusza stylów, sekcji
oraz obiektów, aparatu i pól, które leżą w pliku sąsiednim `studio_postac_obiekty.go`
— jeden kontrakt w jednym miejscu, wzorem `studio.go`. Interfejs
`RepozytoriumPostaciStudia` wchodzi do `RepozytoriumStudia` przez zagnieżdżenie,
bo Studio ma jedno repozytorium, nie dwa, więc adapter modułu dostaje postać tą
samą zależnością, którą dostaje dokument.

Drzewo postaci idzie jednym zapisem, a style i sekcje nie: po drzewie się nie
pyta, drzewo się czyta i zapisuje całe; po stylu i po sekcji się pyta („ile
miejsc używa tego stylu", „która sekcja obejmuje ten znak"), więc mają wiersze.
Wszystkie nazwy pomocnicze tego pliku niosą przedrostek `postac`, bo przestrzeń
nazw pakietu `dane` jest dzielona z innymi plikami modułu.

`PostacJSON` w `PostacDokumentuStudia` niesie drzewo postaci w kształcie
kontraktowego `StudioDocumentForm` — bloki, fragmenty o jednolitej postaci
znaku, tabele, listy.

Blokady fragmentów, dziennik czynności, znakowanie, kopie zapasowe, nastawy
pracy, zajęcia fragmentów i spięcia wykonawców stoją nad tymi samymi tabelami,
ale ich metody deklarują pliki `studio_kontrola_pracy.go` i
`studio_znakowanie_wykonawcy.go`. Interfejs tego pliku ich nie zawiera i nie
zakłada drugich metod nad tymi samymi tabelami: dwa zestawy metod nad jedną
tabelą byłyby dwiema prawdami o tym samym wierszu.

Numer porządkowy postaci przy nadpisaniu rośnie po stronie bazy, nie
wywołującego: gdyby liczył go rdzeń, dwa zapisy z tym samym numerem byłyby
możliwe i okno nie miałoby po czym poznać, że trzyma stan przestarzały.

`PostacDokumentu`: brak wiersza wraca jako `ErrBrakWiersza`, ponieważ dokument
bez zapisanej postaci jest normalnym stanem — dokumenty sprzed wprowadzenia
postaci go nie mają, a rdzeń podstawia wtedy postać domyślną. Zamiana braku na
pustą postać odebrałaby rdzeniowi możliwość odróżnienia stanu „nie ma jeszcze"
od stanu „jest i jest puste".

`StyleDziedziczace` służy jednej rzeczy: zmiana stylu nadrzędnego ma przestawić
wszystkie miejsca, które go używają, także te, które używają go pośrednio przez
styl potomny. Bez tego zapytania rdzeń musiałby czytać cały arkusz i składać
drzewo dziedziczenia przy każdej zmianie.

## budowa/server/internal/dane/przegladarka_karty.go

Trwałość obejmuje tabele `karta_przegladania`, `grupa_kart_przegladania`
i `przestrzen_przegladania`. Skasowanie wiersza karty zamiast oznaczenia jej
kolumną `zamknieta` zabrałoby przestrzeni roboczej to, co zapamiętała: przestrzeń
wskazuje karty, więc karta usunięta z bazy zamieniłaby zapisany zestaw w zestaw
dziurawy, o którym nikt by się nie dowiedział aż do jego otwarcia.

Przynależność do grupy i do przestrzeni stoi po stronie karty. Karta należy do
jednej grupy i jednej przestrzeni naraz, więc skład jednej i drugiej jest
zapytaniem po kolumnie, a nie drugą listą, którą trzeba by prostować przy
każdym zamknięciu karty.

Pole `Przestrzen` w `FiltrKartPrzegladania`: pusta wartość znaczy „wszystkie
karty okna", nie „karty bez przestrzeni". Pole `ZZawieszonymi` istnieje, mimo że
karta zawieszona nadal jest otwarta i domyślnie wchodzi do wykazu, po to, żeby
żądanie mogło wykaz zawęzić do kart żywych. Pole `ZZamknietymi` dopuszcza karty
zamknięte, mimo że rząd kart ich nie pokazuje, bo przestrzeń robocza pamięta
także te, które zostały zamknięte — bez tego przywrócenie zapisanego zestawu
oddawałoby zestaw okrojony i nikt by się o tym nie dowiedział.

Zapytanie wykazu kart zawęża do przestrzeni tym samym idiomem co reszta
wykazów modułu — pusty tekst wyłącza warunek, więc plan zapytania jest jeden,
niezależnie od tego, czy przestrzeń jest podana.

## budowa/server/internal/dane/extension_zaufanie.go

Żaden wiersz tej warstwy niczego nie blokuje: uprawnienie jest zapisem tego, co
manifest deklaruje i co Operator nadał, podpis jest zapisem wyniku weryfikacji,
a referencja sekretu przechowuje klucz jawny, nigdy treść poświadczenia.

Funkcja ZapiszUprawnieniaRozszerzenia wymienia komplet uprawnień jednej strony
zamiast dokładać wiersze do istniejących: manifest, który przestał deklarować
uprawnienie, ma przestać je pokazywać, a nadanie zdjęte przez Operatora ma
zniknąć z tabeli, a nie pozostać w niej jako nieaktualny wiersz.

## przegladarka_wytwory.go

Kolumna odwołania do treści (`tresc_odwolanie`) w tabelach wytworu, zrzutu i
pobrania jest NOT NULL, ponieważ wiersz reprezentuje coś, co istnieje poza
bazą: wytwór i zrzut mają bajty w magazynie plików, pobranie ma plik na
dysku, makro ma kroki możliwe do odtworzenia, a granica obowiązuje kolejne
przebiegi Wykonawcy. Wiersz bez odwołania byłby zapisem bez pokrycia w
rzeczywistym skutku.

## budowa/server/internal/dane/orkiestracja_uklad.go

Plik jest osobny od orchestration.go: tamten prowadzi łuki układu w tabeli
zaleznosc_kroku_automatyki, tutaj mieszkają byty, których łuk nie wyraża.
Bramka mówi, kiedy tory scalają się w jednym kroku, grupa mówi, że zbiór
kroków biegnie razem, a kompensacja mówi, co zrobić, gdy przebieg pękł w pół.
Wszystkie posługują się identyfikatorem zewnętrznym kroku, tak jak łuki: układ
przepisuje się w całości, więc klucz wiersza kroku nie przeżywa przepisania,
a napis Operatora przeżywa.

Kolejki automatyki rozpoznaje nazwa równa jej kodowi, tak zakłada je
core/adapter_modul_automations_kolejka.go i innej drogi nie ma.

Skutek spięcia leży w tabeli kolejka, nie w samym znaczniku spięcia: znacznik
bez przestawienia rodzaju kolejek i wiązania z oknem roli byłby zapisem,
którego nikt nie czyta.

## budowa/server/internal/dane/extension.go

Katalog rozszerzeń stoi poziom wyżej niż ekspert i nie należy do niego. Nie jest
mostem MCP: ten mieszka w tabeli punkt_dostepu, a pozycja katalogu rodzaju mcp
tylko go wskazuje. Nie jest też konektorem eksperta, bo agent_konektor ma
agent_id NOT NULL i należy zawsze do jednego eksperta.

Repozytorium niczego nie pobiera i niczego nie uruchamia. Napis podany w polu
source komendy install ląduje w kolumnie zrodlo_deklarowane jako deklaracja —
zapis faktu, że taki adres podano. Rdzeń pod ten adres nie sięga.

Pole ZmianaRozszerzenia.PunktDostepuID nie ma odpowiednika „odepnij punkt":
kontrakt oznacza accessPointId jako niewymagane w install i configure, ale nie
daje sposobu na jawne odpięcie. Pominięcie pola znaczy więc brak zmiany,
a nie wyczyszczenie — zgadywanie drugiego znaczenia kasowałoby wskazanie mostu
bez żądania.

## isolation.go

Plik nie przechowuje wartości poszczególnych punktów izolacji — jedenaście punktów
izolacji leży w tabeli ustawienie i czyta je repozytorium konfiguracji
(konfiguracja.go, ustawienia_osi.go); ten plik drugiego dostępu do tych samych
wierszy nie zakłada. Tutaj leży profil jako nazwany szablon przełączników,
przypisanie profilu do poziomu oraz warstwa wybrana przez operatora.

Słownik poziomów zasięgu czyta komenda isolation.scope.list: nazwa poziomu i jego
miejsce w kolejności rozstrzygania stoją w bazie i idą stąd do operatora, zamiast
być drugi raz spisane w rdzeniu. Kolejność rozstrzygania nadal należy do pakietu
internal/konfig — tutaj czytany jest opis poziomu, nie reguła.

## budowa/server/internal/dane/design_szablony_materialu.go

Zapis szablonu i jego warstw jest zawsze pełny, wzorem ZapiszKompozycje:
żądanie design.template.save nadsyła komplet warstw, więc repozytorium usuwa
warstwy szablonu i wstawia przysłane od nowa w jednej transakcji. Inaczej
warstwa zdjęta w oknie zostawałaby w bazie i szablon rozchodziłby się z tym,
co Operator widzi. Ten sam rozstrzyg dotyczy stron szablonu: żądanie nadsyła
komplet stron, a dogadywanie różnicy względem stanu zastanego dawałoby dwie
prawdy o tym, które strony publikacja ma.

Baner nie ma stron, a wsteczna zgodność szablonu jednostronicowego nie jest tu
ustępstwem, tylko odzwierciedleniem tego stanu rzeczy.

## library.go

Interfejs RepozytoriumBiblioteki deklaruje wyłącznie ten plik, w całości,
wraz z metodami implementowanymi przez library_wersje.go i
library_kolekcje.go. Interfejs rozdzielony na trzy pliki byłby trzema
prawdami o jednym kontrakcie.

Pole PlikBiblioteki.Sciezka niesie ścieżkę źródłową z wgrania (pole
sourcePath). Bajty treści leżą w magazynie rdzenia i wskazuje je
TrescOdwolanie, więc Sciezka jest zapisem prowenancji, nie drogą do treści.
Pole nie wychodzi jako kontraktowe LibraryFile.path: to pole oznacza drogę
wewnątrz biblioteki, nie ścieżkę systemową maszyny operatora, a klient czyta
pierwszy człon po znaku ukośnika jako katalog nawigacji. Porządek biblioteki
niosą kolekcje — powiązanie wiele-do-wielu, z którego jednej ścieżki nie da
się wyprowadzić — więc pole kontraktu zostaje puste do czasu wprowadzenia
osobnego pojęcia ścieżki w kontrakcie.

Pole SciezkaRepozytorium jest rozłączne ze Sciezka: niesie drogę wewnątrz
biblioteki (kontraktowe LibraryFile.path), którą rozporządza polecenie
przenoszenia pliku.

Kolumna Stan tabeli plik_biblioteki ma warunek CHECK, więc pusty łańcuch
byłby odmową bazy przy każdym wgraniu, które o stanie nie rozstrzyga; funkcja
stanZapisu podstawia wtedy wartość domyślną, a zasób wchodzi do wykazu
czynnego.

Metoda Szukaj realizuje wyszukiwanie po nazwie i po treści jedną frazą
kontraktu, bez rozróżnienia drogi: dopasowanie nazwy działa operatorem LIKE,
dopasowanie treści — indeksem FTS5 zasilanym przy wgraniu pliku i przy
każdej nowej wersji. Plik niezaindeksowany albo o treści nietekstowej nie
znika z wyszukiwania, ponieważ dopasowanie nazwy jest osobnym członem
alternatywy, nie warunkiem dodatkowym zależnym od obecności wiersza w
indeksie.

## budowa/server/internal/transport/serwer_test.go

Sprawdziany warstwy nasłuchu idą przez prawdziwe gniazdo, nie przez atrapę
biblioteki, ponieważ wszystko, co w tej warstwie potrafi zawieść, zawodzi na
styku dwóch pętli, gniazda i rejestru połączeń — dokładnie tam, gdzie atrapa
niczego by nie odwzorowała. Serwer wstaje na porcie wskazanym przez system,
więc sprawdziany idą równolegle i nie biją się o port.

W sprawdzianie TestRdzenNiepodlaczonyOddajeOdmoweZamiastCiszy typ odpowiedzi
jest typem komendy, nie zdarzeniem *.unknown: komenda kontraktowa została
rozpoznana przez rejestr transportu, więc odmowa wraca pod jej własną nazwą.
Zdarzenie obszaru dostaje wyłącznie typ spoza kontraktu. Dla klienta rozstrzyga
i tak stan wraz z kodem, po których wie, że rdzeń uchwytu nie ma.

## automatyka_petla.go

Definicja, harmonogram i przebieg opisują automatykę; trzy byty tego pliku opisują
jej obieg: kto bierze udział w pętli, który agent prowadzi krok i na co bieg czeka.

Rozszerzenie idzie osobnym interfejsem RepozytoriumPetli, po który rdzeń sięga
asercją typu na porcie automatyk — tak jak po Harmonogramy i Orkiestracja
w pliku kompozycja.go.

Bieg oczekujący czeka na sygnał z zewnątrz, nie na zegar ani na operatora, więc jego
stan musi przeżyć restart rdzenia — stąd wiersz w bazie, a nie wpis w mapie adaptera.

Tabela obsada_biegu niesie ten sam byt dla biegu automatyki i biegu orkiestracji,
którego wiersza w tabeli automatyka nie ma.

ZapiszObsade zastępuje obsadę w całości zamiast dopisywać wiersze, ponieważ obsada
jest wykazem zamkniętym i scalanie zostawiałoby uczestników, których operator z niej
usunął. Obsada pusta jest stanem poprawnym — automatyka bez obsady biegnie na modelu
wskazanym w kroku.

## design_makiety.go

Warstwa wchodzi do kompozycji pojedynczo, a nie kompletem, w odróżnieniu od
obszaru kompozycji, który zna wyłącznie zapis pełny (usunięcie całości i
wstawienie kompletu, bo zapis planszy nadsyła zawsze komplet). Makieta
pracuje inaczej: układ automatyczny przestawia położenia warstw już
leżących, a instancja komponentu dokłada jedną warstwę do planszy, na której
stoją inne. Przepisywanie całej planszy przy każdym takim ruchu kasowałoby
warstwy, o których wołający w danym żądaniu nic nie mówił. Stąd dwie osobne
czynności: PrzestawWarstweKompozycjiDesignu i DolozWarstweKompozycjiDesignu.

Przynależność warstwy do ramki oraz jej więzy responsywne wiszą na
identyfikatorze zewnętrznym warstwy, nie na kluczu wiersza. Powód: zapis
planszy przepisuje wiersze warstw od nowa, więc klucz wiersza warstwy nie
przeżywa zwykłego zapisu kompozycji, a identyfikator zewnętrzny przeżywa,
ponieważ klient nadsyła go z powrotem przy każdym żądaniu.

## budowa/server/internal/models/adapter_obrazy.go

kanalObrazow jest bliźniakiem KanalAPI i dzieli z nim całą warstwę dostępową:
nagłówki, klucz z sejfu albo ze zmiennej środowiskowej, dodatki ciała, limit
czasu — jedna prawda o poświadczeniach, jedna o nagłówkach. Różni się kształtem
żądania (prompt zamiast wiadomości) i tym, że wynik nadaje fragmentem image,
nie porcjami tekstu.

Kształt żądania jest ten, który dostawcy powtarzają za OpenAI Images: POST
z ciałem {model, prompt, n, size} i odpowiedzią {"data":[{"b64_json"|"url"}]}.
Powtarzają go dziś także dostawcy niezależni i bramy lokalne, więc jest to
najczęstszy kształt, a nie jeden dostawca zaszyty w kodzie. Parametry wiersza
rejestru: base_url (pełny adres punktu końcowego, wymagany), rozmiar (pole
size, domyślnie 1024x1024), liczba (pole n, domyślnie 1), format_odpowiedzi
(pole response_format, pusty = nie wysyłamy), sciezka_base64 (domyślnie
data.0.b64_json), sciezka_adresu (domyślnie data.0.url), typ_tresci
(domyślnie image/png), a także sciezka_bledu, naglowki, naglowek_klucza,
przedrostek_klucza, cialo_dodatkowe i limit_sekund, jak w kanale api.

Wiersz zakłada się channel.add z rodzajem api oraz parametrem adapter równym
"obrazy". Rodzaj mówi, jak kanał rozmawia (HTTP), adapter mówi, co oddaje,
dlatego rodzaju kanału nie trzeba dokładać ani migracją, ani do kontraktu.

Wywołanie niosące Zapytanie.ObrazyWejsciowe jest edycją materiału, nie
generowaniem od zera, i dostawcy dają na nią osobny punkt końcowy o osobnym
kształcie żądania (u OpenAI images/edits: multipart/form-data z plikiem image
i opcjonalną maską mask). Bramy lokalne częściej przyjmują ten sam adres
z bajtami w polu JSON. Adapter zna oba kształty, bo obu nie da się pogodzić,
a zgadywanie kończyłoby się odmową dostawcy zamiast obrazu. Parametry:
adres_edycji (pusty = ten sam co base_url), postac_obrazu (wieloczesciowa
domyślnie albo base64), pole_obrazu (domyślnie image), pole_maski (domyślnie
mask). Wywołanie bez obrazu wejściowego idzie tą samą drogą co dotąd: pole
puste nie zmienia ani adresu, ani kształtu ciała.

Wiersz bez poświadczenia buduje się celowo: odmowa ma paść w chwili wywołania,
wymieniając brak z nazwy, a nie zniknąć jako kanał pominięty w wykazie
rejestru. Kanał tekstowy dopuszcza brak odwołania, bo punkt końcowy bez
uwierzytelnienia istnieje (Ollama pod adresem pętli zwrotnej); punkt końcowy
generujący obrazy bez klucza nie istnieje, więc brak odwołania rozstrzyga się
przed żądaniem, zamiast wracać jako cudzy błąd dostawcy. To samo dotyczy
odwołania, które jest, ale nie prowadzi do sekretu — sejf bez wpisu, zmienna
środowiskowa nieustawiona: różnica jest widoczna w treści komunikatu, bo
Operator poprawia te dwa stany w różnych miejscach.

Materiał nieczytelny w cialoWieloczesciowe jest odmową, nie wysyłką bez
materiału: żądanie samego polecenia wróciłoby obrazem wygenerowanym od zera,
a Operator prosił o obróbkę swojego zdjęcia i dostałby cudze bez ani jednego
słowa o podmianie.

## budowa/server/internal/dane/extension_cykl.go

Interfejs RepozytoriumRozszerzen deklaruje plik extension.go; ten plik i trzy
sąsiednie — extension_protokol.go, extension_integracje.go, extension_zaufanie.go
— dokładają mu metody. Jedno repozytorium, cztery pliki wedle odpowiedzialności,
tak jak w obszarze Apps.

## kondycja.go

Repozytorium nie wykonuje pomiaru i nie zna żadnego rodzaju sondy: zapisuje definicję
i zapisuje wynik, który ktoś zmierzył. Rozdział jest tu istotny — gdyby warstwa danych
umiała ustawić stan sondy wprost, istniałaby droga do odłożenia stanu bez pomiaru,
a właśnie tego rodzina health zabrania.

Zapis wyniku podnosi jednocześnie odbicie w definicji (pola ostatni_stan
i ostatni_przebieg) w jednej transakcji, żeby wykaz sond nie pokazywał stanu innego
niż ostatni wiersz serii.

## trwalosc_niszczaca_test.go

Rdzeń wykonuje trzy drogi kasujące dane przy każdym starcie: kaskada
schematu zabiera okna i wiadomości razem z sesją, kosz kasuje trwale po
terminie, a retencja przycina historię sama. Każda z nich jest
nieodwracalna i żadna nie pyta Operatora. Sprawdzian mierzy dwie rzeczy
naraz: że każda droga kasuje to, co ma kasować, i że żadna nie tyka niczego
poza tym. Druga część jest ważniejsza, ponieważ nadmiarowe skasowanie nie
zgłasza się błędem, tylko brakiem danych, który łatwo przeoczyć.

## budowa/server/internal/konfig/rozgloszenie_test.go

Bez drogi na żywo Operator zmienia nastawę i nic się nie dzieje, dopóki czegoś
nie przeładuje. Z nią dzieje się dokładnie tyle, ile trzeba: zapis na poziomie
szerszym nie budzi nikogo, kto ma wartość z węższego. Ta druga połowa jest
ważniejsza od pierwszej, bo doręczenie nadmiarowe nie wygląda na błąd —
wygląda na odświeżenie.

## budowa/server/internal/dane/tlumaczenie_kontrola.go

Wykaz terminów zawężony, profile kontroli jakości, obieg zatwierdzeń panelu
i ustalenia korekty językowej mieszkają w jednym pliku, bo wszystkie trzy
odpowiadają na jedno pytanie: czy ten przekład wolno wypuścić. Rozbicie ich
na osobne pliki dałoby trzy nagłówki mówiące to samo.

ZapiszProfilQa zakłada profil albo nadpisuje zastany i wymienia jego kontrole
w całości, bo kontrakt nadsyła wykaz kontroli kompletem.

ZapiszZatwierdzenie dokłada krok obiegu i przestawia migawkę panelu w jednej
transakcji; inaczej wiersz obiegu i stan panelu rozjechałyby się przy awarii
między dwoma zapisami.

ZapiszUstaleniaKorekty wymienia otwarte ustalenia panelu na nadesłane.
Ustalenia rozstrzygnięte, zastosowane albo odrzucone, zostają, bo kolejny
przebieg nie ma prawa skasować odpowiedzi Operatora.

RozstrzygnijUstalenieKorekty znakuje ustalenie jako zastosowane albo
odrzucone, zapisując chwilę, nie wartość logiczną: kiedy niesie więcej niż
czy, a czy da się z kiedy odczytać.

## roundtable_graf.go

Graf jest trwały, ponieważ oznaczenie węzła jako kluczowego stawia operator poleceniem
roundtable.argument.pin. Oznaczenie postawione na węźle wyliczanym w locie znikałoby
przy następnym odczycie razem z identyfikatorem węzła.

Ponowna analiza zastępuje graf tury, a nie dokłada się do niego: dwa przebiegi
wydobycia argumentów na tym samym zapisie dałyby każdy węzeł dwa razy. Funkcja
ZastapGrafDebaty przenosi oznaczenia kluczowe po treści węzła, żeby zastąpienie grafu
nie zgubiło wyboru operatora.

## budowa/server/internal/dane/design_wektor.go

Zapis ścieżki jest zawsze pełny. Kontrakt DesignVectorPathSetRequest.Nodes
nadsyła komplet węzłów przy każdej zmianie — zmiana jednego węzła idzie tą samą
drogą co narysowanie ścieżki od zera. Repozytorium nie dogaduje różnicy
względem kształtu zastanego: zna wyłącznie „załóż albo nadpisz" po
identyfikatorze zewnętrznym, wzorem ZapiszZestawZetonowDesignu.

Węzły, wypełnienie i obrys jadą zapisem JSON w kolumnie. Powód stoi w nagłówku
migracji 316: żadne zapytanie nie pyta o pojedynczy węzeł, a wiersz na węzeł
znaczyłby kasowanie i wstawianie kilkudziesięciu wierszy przy każdym drgnięciu
pióra. Warstwa danych nie zagląda w treść tego zapisu — składa go i rozkłada
adapter, bo to on zna kontrakt.

Symbol i jego członkowie: definicja symbolu jest wykazem ścieżek i warstw,
który podmienia się w całości. Liczba członków jest tym, co rdzeń oddaje jako
liczbę miejsc, do których zmiana doszła — liczbą wierszy naprawdę zapisanych,
nie obietnicą.

## aplikacje_srodowiska.go

Zmienna środowiskowa niesie wartość jawną albo odwołanie do sekretu, nigdy obie naraz.
Warunek CHECK schematu pilnuje tego po raz drugi, a warstwa dane odrzuca obie wartości
podane jednocześnie, żeby literówka wołającego wracała czytelnym powodem, a nie treścią
błędu SQL.

## prowenancja.go

Ślad jest podstawą dwóch rzeczy naraz: Provenance Explorer pyta o pojedyncze
wywołanie, a rozliczenie zużycia liczy sumy po wymiarze. Obie odpowiedzi powstają
z tej samej tabeli — druga tabela z tymi samymi liczbami rozjechałaby się
z pierwszą przy pierwszej korekcie cennika.

Zawężenie wykazu idzie do bazy, nie do pętli w Go. Wywołań przybywa w tempie pracy
operatora, a wykaz bez zawężenia po czasie potrafi zająć całą tabelę.

Odwzorowanie wymiaru rozliczenia na kolumnę tabeli stoi w repozytorium, a nie
w adapterze, ponieważ to warstwa danych wie, którą kolumną grupuje. Wymiar spoza
wykazu jest odmową, nie cichym zejściem do wymiaru domyślnego — suma policzona
po innej osi niż zamówiona wyglądałaby identycznie jak zamówiona i nie dałoby się
ich odróżnić.

## design.go

Interfejs RepozytoriumDesignu deklaruje wyłącznie ten plik, w całości —
także metody obszarów zasobów i kompozycji — żeby kontrakt stał w jednym
miejscu; implementację niosą design_zasoby.go i design_kompozycje.go. Prompt
ma własną tabelę, ponieważ wiele zasobów i wariantów powstaje z tego samego
promptu: prompt_design jest jedną prawdą o promptcie, a
zasob_design.prompt_id jedynym odwołaniem.

Metoda UstawUlubionyZasobu stoi osobno od ZapiszZasob z podmienionym polem,
ponieważ pełny zapis wymaga kompletu pól zasobu: oznaczenie ulubionego przez
odczyt-zmień-zapisz nadpisywałoby przy wyścigu dwóch komend cudzą zmianę
nazwy czy wymiarów wartością sprzed odczytu.

Ikony własne: w bazie leżą wyłącznie ikony własne, ponieważ katalog ikon
otwartoźródłowych jest wkompilowany w binarium rdzenia — repozytorium go nie
zna i nie ma go czym zasiać ani zgubić.

Warsztat fotografii: łańcuch edycji nie dubluje wariantów zasobu, bo że
wariant powstał ze źródła, mówi kolumna zasob_design.wariant_zasobu_id.
Metody warsztatu trzymają to, czego ta kolumna nie niesie — czynność, jej
nastawy i drogę rachunku.

Ścieżki wektorowe: ścieżka jest bytem osobnym od warstwy, ponieważ krzywa
nie mieści się w prostokącie warstwy, a operacja logiczna potrzebuje obu
krzywych z osobna.

Ramki makiety: przynależność warstwy do ramki wiąże się identyfikatorem
zewnętrznym warstwy, ponieważ zapis planszy przepisuje komplet warstw od
nowa i klucz wiersza warstwy tego zapisu nie przeżywa.

Gradienty wypełnienia: gradient rozstrzyga się celem — kompozycją, ścieżką
albo warstwą — nie identyfikatorem osobnym: zapis gradientu zakłada albo
zmienia gradient stojący na wskazanym bycie, a nie dokłada drugi obok.

Funkcja noweRepozytoriumDesignu przyjmuje parametr bazy danych osobno,
ponieważ obszar kompozycji prowadzi zapis pełnej listy warstw w transakcji
(usuń i wstaw od nowa), której zapis promptu nie potrzebuje — mieści się w
jednym poleceniu.

## budowa/server/internal/konfig/rozstrzyganie_test.go

Rozstrzygacz zasięgu jest jedyną logiką pakietu, w której błąd nie wywraca
niczego. Ustawienie rozstrzygnięte o jeden poziom za szeroko przecieka między
sesjami Operatora, a widać to dopiero po fakcie — po zachowaniu modelu, nie po
komunikacie. Reguła ma dwa piętra: POZIOM rozstrzyga pierwszy, OŚ dopiero
w ramach poziomu. Odwrócenie tego znaczyłoby, że wybór konta unieważnia decyzję
podjętą wprost w oknie komunikacji — czyli odbiera Operatorowi sterowanie
zamiast je rozszerzać. Kompilacja tej reguły nie pilnuje, bo obie drogi
zwracają wartość tego samego typu.

TestNajwezszyZapisWygrywa jest sprawdzianem, dla którego ten plik powstał.
Reguła schodzenia w górę jest jedyną drogą, którą Operator odzyskuje ustawienie
szersze po skasowaniu węższego — pomyłka na którymkolwiek szczeblu zostawia go
z wartością, której nigdzie nie zapisał.

## aplikacje.go

Interfejs RepozytoriumAplikacji stoi w całości w tym pliku, wraz z metodami,
które implementują pozostałe pliki obszaru Apps. Interfejs rozdzielony na
trzy pliki byłby trzema prawdami o jednym kontrakcie.

Pole ArchitekturaApp.RoznicaWersji odkłada je metoda ZapiszArchitekture w
wierszu tabeli wersja_architektury_apps. Odczyt architektury tego pola nie
wypełnia, ponieważ historia ma własny odczyt WersjeArchitekturyApp, a pole
wypełniane w obie strony sugerowałoby, że wiersz bieżący pamięta, czym
różnił się od poprzednika.

Metoda ArchitekturaOkna czyta po oknie, ponieważ klient po odświeżeniu zna
wyłącznie okno, a kodu architektury nadanego przy pierwszym zapisie już nie
pamięta. Indeks na oknie i znaczniku czasu malejąco sprawia, że odczyt
najświeższej definicji nie skanuje tabeli.

Metoda PlikWarsztatu zwraca jeden plik po kluczu naturalnym, ponieważ zapis
zmiany musi powiedzieć, czy plik powstał, czy został zmieniony, a operacja
UPSERT sama tego nie mówi — obie kolumny czasu mają osobne wartości
domyślne, więc porównanie znaczników byłoby zgadywaniem.

Metoda Wdrozenia liczy wszystkie wdrożenia spełniające warunki osobnym
zapytaniem COUNT, nie długością zwróconej strony — inaczej łączna liczba
przy limicie mniejszym niż dziennik kłamałaby o rozmiarze historii.

Metoda ZapiszArchitekture nadsyła całą listę komponentów i zależności na
nowo, bez trybu częściowej zmiany — zapis jest więc zawsze usunięciem
zastanych wierszy i wstawieniem od nowa. Wersja rośnie przy każdym zapisie
definicji; operacja UPSERT ustawia to w klauzuli ON CONFLICT. Numer wersji
zapisywany w historii czytany jest z wiersza po UPSERT-cie, ponieważ to on
go podniósł — wartość policzona osobno rozjechałaby się z bazą przy dwóch
zapisach naraz.

## stan.go

Cały pakiet stoi na bibliotece go-git wkompilowanej w binarium rdzenia. Nie ma tu ani
jednego uruchomienia procesu i mieć nie będzie: funkcja zależna od programu, którego
instalka nie niesie, jest u odbiorcy odmową, a nie funkcją. Odczyt stanu repozytorium
ma działać zawsze, bo od niego zaczyna się każda inna czynność okna Git Panel.

Pakiet nie zna kontraktu komend, zna wyłącznie kształt danych. Stąd nie wychodzi ani
jedna odmowa protokołu: pakiet oddaje wynik albo błąd Go, a nazwanie go operatorowi
należy do adaptera. Dzięki temu ten sam silnik obsługuje Git Panel, margines zmian
w edytorze i przegląd różnicy przed scaleniem, nie ucząc się o żadnym z nich.

## budowa/server/internal/dane/badania_raport.go

Ten plik jest trzecią częścią RepozytoriumBadan zadeklarowanego w badania.go;
źródła leżą tam, ustalenia w badania_ustalenia.go.

Zapis raportu jest zawsze pełny, po wzorze ZapiszKompozycje: komplet sekcji
jest usuwany i wstawiany od nowa w jednej transakcji, bo research.report.build
nie zna trybu częściowej zmiany. Przestrzeń badania jest jednowierszowa, bo
kontrakt research.workspace.set nie niesie identyfikatora — UstawPrzestrzen
nadpisuje jedyny wiersz o id = 1. Tabela pusta w Przestrzen wraca jako zakres
i lista puste, nie jako błąd: to stan startowy.

Kolumna Cel doszła migracją 150 — do niej eksport zawsze szedł do pobrania,
bo innego celu kontrakt wtedy nie miał.

## budowa/server/internal/dane/zestaw.go

Zestaw składa pod jednym dachem repozytoria każdego obszaru danych aplikacji.
Kilka pól opisuje dwa widoki tej samej tabeli albo implementacji, celowo
oddzielone, bo odpowiadają na różne pytania: KonfiguracjaOsi jest Konfiguracją
widzianą razem z osią rozstrzygania; Macierz czyta wiersz tabeli
srodowisko_modul w całości, podczas gdy Moduly zna ją wyłącznie jako filtr
wykazu; WersjeAgenta i ArchiwumAgentow to jedna implementacja widziana jako
historia tożsamości eksperta i jako archiwum; Historia pyta o wiersz
wiadomosc z drugiej strony niż Wiadomosci — po identyfikatorze kontraktowym
okna, od najnowszej, kursorem czasu — i jako jedyne kasuje; KoszSesji jest
odwrotną stroną tabeli sesja, widzi wyłącznie wiersze ze znacznikiem
usunieto_o, których wykaz sesji żywych nie widzi wcale.

ZakresAgenta odpowiada nie na pytanie, jaki jest ekspert, lecz co temu
ekspertowi wolno zrobić w systemie: moduły zastosowania, zakresy izolacji
technicznej, granica Subagent Network, wpisy uprawnień oraz konektory
i przypisania od strony eksperta — dlatego jest osobnym repozytorium od
biblioteki ekspertów.

Przekazania nie powtarzają więzi koordynator-wykonawca: ta mieszka w kolumnie
okno_komunikacji.okno_koordynatora_id.

WylaczeniaPamieci stoi osobno od Pamiec i PrzestrzenRobocza, bo wyłączenie nie
jest wpisem pamięci ani jego zmianą — nie dotyka treści i znosi się jednym
ruchem. WyciszeniaNakladki jest bytem rdzenia, nie stanem jednego okna: bez
wspólnego wiersza Operator wyciszałby w jednej powłoce, a w drugiej sugestie
wchodziłyby dalej.

Zestaw pól baza w strukturze Zestaw obsługuje repozytoria składane na żądanie
poza konstruktorem Otworz (Rozszerzenia, NarzedziaSesji), które potrzebują
bezpośredniego połączenia do własnych transakcji.

## studio_postac_obiekty.go

Aparat dokumentu i pola stoją w jednym pliku, nie osobno, ponieważ pracuje
się nimi tą samą drogą: oba są przypięte do miejsca w treści, oba bywają
nieświeże i oba odświeża się wykazem, nie po jednym wierszu. Rozdział na
osobne pliki byłby rozdziałem na papierze — kod odczytu i zapisu byłby ten
sam dwa razy. Nazwy pomocnicze tego pliku niosą przedrostek postac.

## budowa/server/internal/dane/extension_protokol.go

Metryka użycia liczy się z wierszy, nie z licznika: liczba wywołań, liczba
niepowodzeń i średni czas w oknie czasu wymagają trzech różnych agregatów nad
tym samym zbiorem.

ZapiszNarzedziaRozszerzenia wymienia komplet wpisów odkrytych u integracji
zamiast dokładać nowe: serwer, który przestał udostępniać narzędzie, ma
przestać je pokazywać, a wpis pozostawiony byłby obietnicą bez pokrycia.

## orkiestracja_podagenci.go

Podagent to zadanie w tle, którego trwałą tożsamością jest pozycja kolejki, a proces
modelu jest wyłącznie sposobem jej wykonania. Stąd kształt tego repozytorium: wiersz
zna swoją pozycję kolejki (pole pozycja_kolejki_id), a cyklu życia zlecenia nie
prowadzi — prowadzi go silnik kolejek.

Nie ma tu odczytu tabeli pozycja_kolejki: jej czytelnikiem jest repozytorium kolejek.
Wiersz podagenta niesie wyłącznie to, czego pozycja nie wie — pod kim biegnie, jak się
nazywa i ile kosztował.

Repozytorium wchodzi metodą zestawu, nie polem — ten sam wzorzec co
Zestaw.Rozszerzenia() i Zestaw.RoleOkien(): rejestr nie trzyma stanu poza wskaźnikiem
na wspólną pamięć zapytań, więc złożenie go na żądanie kosztuje tyle, co odczyt pola.

Znaczniki czasu w poleceniu ustawStanPodagenta stawia baza, nie rdzeń: chwila
rozpoczęcia zapisuje się raz, przy pierwszym wejściu w bieg, a chwila zakończenia raz,
przy pierwszym stanie końcowym. Dzięki temu powtórzony zapis stanu nie przesuwa
historii podagenta. Powód zakończenia i oznaka życia dokładają się tym samym
poleceniem: stan mówi co, powód mówi dlaczego, a oznaka — kiedy rdzeń ostatni raz tego
wiersza dotknął. Powód dla stanu niekońcowego zostaje zastany, bo przejście ze stanu
oczekującego do biegnącego niczego nie kończy.

## alerty.go

Repozytorium nie ewaluuje reguł i nie zna ani jednej miary. Trzyma
definicję oraz zapis wyzwolenia wraz z wartością zmierzoną w chwili
wyzwolenia. Ewaluacja należy do adaptera, ponieważ miary pochodzą z
magazynów, o których warstwa danych alertu nie ma prawa wiedzieć: ślad
wywołań, dziennik błędów, seria pomiarów sond.

Metody LiczbaNieudanychPomiarowSond i LiczbaNieudanychPozycjiKolejki niosą
dwie miary, których nie da się wziąć z magazynu śladu wywołań ani z
dziennika błędów — a bez nich reguły miar probeFailure i processFailure
nigdy by się nie wyzwoliły. Zapytania stoją w tym repozytorium, ponieważ to
ewaluacja reguły ich potrzebuje; żadne inne repozytorium ich nie woła.

## budowa/server/internal/dane/terminal_wyposazenie_zapis.go

Usunięcie oddaje prawdę o skutku (bool), a nie samo „nie było błędu". Kontrakt
komend terminal.host.remove, terminal.script.remove i terminal.key.remove
mówi wprost: fałsz znaczy, że wpisu nie było, i nie jest błędem. Bez policzenia
zmienionych wierszy rdzeń nie miałby czym tego rozróżnić i odpowiadałby prawdą
zawsze.

W ZapiszSkrypt odczytanie numeru wersji osobnym zapytaniem przed zapisem
dałoby dwóm równoległym zapisom ten sam numer, a warunek UNIQUE na parze
pozycja-wersja odrzuciłby drugi z nich; numer nadaje więc baza wyrażeniem
wersja + 1 wykonanym w tej samej transakcji co wpis wersji.

## budowa/server/internal/dane/zlecenia_kolejki.go

Tabela zlecenie_kolejki nie jest drugim silnikiem kolejek ani drugą tabelą
pozycji obok pozycja_kolejki: ta ostatnia opisuje etap pętli
koordynator-wykonawca (tytuł, treść zlecenia, werdykt weryfikacji, licznik
obiegów), podczas gdy kontraktowy QueueItem niesie ładunek strukturalny,
priorytet, termin wykonania, klucz idempotencji i warunek przetworzenia.
Wtłoczenie jednego w drugie kazałoby kolumnie tytul nieść ładunek,
a werdykt_weryfikacji stan o zupełnie innym słowniku.

Zlecenie martwe zostaje w tej samej tabeli, ze stanem martwe. Kolejka zadań
martwych jest widokiem, nie osobnym magazynem: queue.dead.list bez wskazania
kolejki oddaje zadania martwe wszystkich kolejek, więc przeniesienie ich
gdzie indziej odebrałoby im pochodzenie.

Stan pusty w wykazie zleceń kolejki znaczy wszystkie stany: zadania martwe
i zdjęte wychodzą wtedy razem z resztą, bo Queue Manager pokazuje je
w kolumnie stanu, a nie ukrywa przed Operatorem.

PolitykaKolejki: kolejka bez zapisanej polityki nie jest kolejką bez
polityki, bo kolumny mają wartości domyślne modelu konfiguracji.

ZlecenieKluczem: ten sam klucz idempotencji w dwóch kolejkach opisuje dwa
różne zlecenia dwóch różnych torów, więc odczyt zawsze pyta o parę
kolejka-klucz, nie o sam klucz.

ZleceniaKolejki: wykaz zleceń bywa przycięty granicą wyniku, a licznik
wszystkich zleceń kolejki nie, dlatego funkcja zwraca oba osobno.

Odcinek głębokości kolejki wraca z zapytania jako numer, więc chwilę
odtwarza się mnożeniem numeru przez długość odcinka — tak powstaje początek
odcinka w sekundach epoki.

## aplikacje_wytwory.go

Cztery byty tego pliku — podgląd, motyw, dziennik i artefakty — łączy to, że
każdy wskazuje coś poza bazą: adres stojącego serwera, wiersz dziennika
wytworzony przez pracę, plik archiwum w magazynie treści rdzenia. Wiersz bez
tego czegoś byłby meldunkiem bez skutku, dlatego zapisuje go wyłącznie kod,
który ten skutek właśnie wywołał.

## budowa/server/internal/dane/kroki_zlecenia.go

Sterowanie krokiem oddaje wyłącznie to, czego pozycja_kolejki nie umie
powiedzieć: że krok stoi, na jaką decyzję czeka, jaka decyzja zapadła i czy
dojechała do wykonawcy. Metody siedzą na repozytoriumKolejek, a nie na
własnym typie, bo wstrzymanie kroku i stan kroku to dwa pytania o tę samą
tabelę pozycja_kolejki i ten sam dziennik log_akcji_kolejki: osobne
repozytorium musiałoby powtórzyć odczyt położenia pozycji i własnym zapisem
dziennika rozjechać się z zapisem silnika, wzorem par Moduly/Macierz
i Sesje/KoszSesji — jedna implementacja, dwa widoki. Rozszerzenie widać przez
RepozytoriumWstrzymanKroku: port sięga po nie asercją typu, tak jak rdzeń
sięga po wiazaneKolejki przy queue.list i queue.link.

Stany czeka i biegnie nie należą do słownika kolumny wstrzymanie_kroku.stan:
czyta się je z pozycja_kolejki.stan, a zapisane drugi raz byłyby drugą
prawdą.

Akcje dziennika kolejki zapisywane przez sterowanie krokiem idą do tego
samego log_akcji_kolejki, co działania silnika — jeden ślad, nie dwa.

CzyZastosowane: epizod zdecydowany, lecz niezastosowany, to decyzja w
drodze — nie wolno jej zgubić.

RepozytoriumWstrzymanKroku wypełnia ta sama implementacja, co
RepozytoriumKolejek — drugiej nie ma.

Zapis decyzji wchodzi wyłącznie na epizod jeszcze nierozstrzygnięty
i jeszcze niezastosowany. Warunek stoi w SQL, nie w warstwie wyżej, bo
inaczej dwa równoległe rozstrzygnięcia nadpisałyby się nawzajem. Zastosowanie
wchodzi wyłącznie na epizod rozstrzygnięty i jeszcze niezastosowany, więc
powtórne wywołanie nie zrobi drugiego zastosowania.

WstrzymajKrok: krok w stanie końcowym odbija się o CHECK schematu, a odmowa
jest wtedy prawdą o kroku, nie awarią zapisu.

CzynneWstrzymanie: brak epizodu nie jest błędem, bo większość kroków nigdy
nie była wstrzymana.

ZapiszDecyzje: epizod już rozstrzygnięty albo już zastosowany nie zostaje
ruszony — zapytanie ma warunek i brak trafienia jest tu odmową, nie ciszą.

OznaczZastosowanie: doręczenie zapisuje się osobno od zastosowania, bo to
dwa różne fakty — krok ruszył zgodnie z decyzją, i decyzja dojechała do
tego, kto pracę wykonuje.

CzynneWstrzymaniaKolejki: indeks częściowy schematu gwarantuje najwyżej
jeden epizod na krok, więc mapa po identyfikatorze pozycji niczego nie
gubi.

HistoriaWstrzymanKroku: wiersze zastosowane zostają w wykazie, bo ślad
wstrzymania jest częścią przejrzystości pętli.

pozycjaWstrzymania odczytuje krok, którego dotyczy epizod, w tej samej
transakcji co zapis, żeby dziennik nie wskazał innego kroku niż zapis.

## budowa/server/internal/dane/aplikacje_produkt.go

Interfejs RepozytoriumAplikacji oraz typ *repozytoriumAplikacji deklaruje
aplikacje.go; ten plik dokłada mu metody, nie drugi kontrakt.

Produkt jest jeden na okno, więc zapis jest UPSERT-em po kolumnie okno,
rozstrzygnięcie zapisane na czole migracji: kod zewnętrzny nadany przy
pierwszym zapisie zostaje, bo AppProduct.id ma być stały, a kolejne zapisy
zmieniają metadane, nie tożsamość produktu. Etap i kamień milowy mają własne
identyfikatory zewnętrzne, bo obie komendy zapisu potrafią wskazać byt
zmieniany — zapis jest tam UPSERT-em po identyfikatorze.

KamienMilowyApp niesie kody etapów, które się na niego składają, bo kontrakt
oddaje kamień milowy zawsze razem z nimi — rozdzielanie tego na dwa odczyty
byłoby pracą bez odbiorcy. ZapiszKamienMilowyApp wymienia ten związek
w całości w jednej transakcji, bo kontrakt nadsyła stageIds bez trybu
częściowej zmiany.

UsunKamienMilowyApp zwraca prawdę, gdy wiersz naprawdę zniknął: kontrakt
oddaje pole deleted, więc "nie było czego kasować" nie ma prawa wrócić jako
powodzenie usunięcia.

etapyKamieniOkna czyta jednym zapytaniem cały wykaz etapów okna zamiast
jednego zapytania na każdy kamień milowy, bo wykaz kamieni ciągnąłby inaczej
tyle zapytań, ile ma pozycji.

## studio_katalogi.go

Operacje własne Tools Panel, łańcuchy operacji i profile wydania mieszkają
w zasięgu konfiguracji (globalny, środowisko, projekt, sesja), a nie przy
dokumencie: Operator zapisuje własny prompt raz i sięga po niego w każdym
dokumencie. Wiązanie ich z dokumentem kazałoby przepisywać je przy każdym
nowym pliku.

Gałąź i odwołanie do wersji należą do dokumentu, ale stoją w tym pliku,
ponieważ są bytami rejestru — mają własny cykl życia i własny identyfikator
zewnętrzny, inaczej niż zmiana śledzona, która bez dokumentu nie znaczy nic.

Kontrakt repozytorium mówi nazwami dziedziny, na przykład zapisz operację, a
nie nazwą tabeli podanej parametrem: wołający nie ma rozstrzygać, w której
tabeli byt mieszka. Opis tabeli zostaje szczegółem tego pliku.

## budowa/server/internal/repozytorium/skutek_repozytorium_test.go

Materiał każdego sprawdzianu powstaje w samym pliku: repozytorium zakładane na czas
sprawdzianu, z zatwierdzeniami wytworzonymi tą samą biblioteką. Repozytorium wniesione do
drzewa zestarzałoby się razem z wersją biblioteki, a sprawdzian oparty na repozytorium
produktu mierzyłby przypadek, ponieważ jego stan zmienia się przy każdej pracy. Żaden
sprawdzian nie pomija się przy braku programu git w środowisku uruchomieniowym: pakiet nie
uruchamia ani jednego procesu zewnętrznego.

## biblioteka_reguly.go

Kolekcje mają tu drugie wejście obok `library_kolekcje.go` i to nie jest powielenie: tamten plik
odpowiada za założenie kolekcji i przypisanie zasobów, ten za odczyt kolekcji jako bytu opisanego
(rodzic, reguła, licznik). Podział idzie po pytaniu, nie po tabeli.

Reguła nie wykonuje się sama. Warstwa danych przechowuje warunek i wskazuje kolekcję docelową;
przeliczenie, czyli zamiana warunku na wykaz zasobów, należy do rdzenia, bo to on zna znaczenie
członów warunku.

## workspace_zadania.go

Wykaz zadań schodzi z bazy w całości dla jednego projektu, a zawężenia (stan, wykonawca, etykieta,
fraza, termin) rozstrzyga rdzeń. Powód jest jeden: zawężeń jest siedem i wchodzą w dowolnym
połączeniu, więc zapytanie składane z kawałków byłoby napisem budowanym w locie, a projekt liczy
zadania w setkach, nie w milionach.

## budowa/server/internal/dane/historia.go

Repozytorium historii nie zakłada własnej tabeli: pozycją historii jest wiersz tabeli
wiadomości, tylko on niesie rolę, treść i czas. Osobne repozytorium obok repozytorium
wiadomości bierze się z drugiego pytania o tę samą tabelę — tamto prowadzi turę, pisze
wypowiedź i czyta okno po kluczu wewnętrznym, nigdy nie usuwa; historia pyta po
identyfikatorze kontraktowym okna, od najnowszej, kursorem czasu, i jako jedyna kasuje.
Usunięcie zabiera też bloki wiadomości, ponieważ bez klucza obcego kaskada bazy ich nie
sprząta.

Liczba pozycji całego okna pomija warunek kursora celowo: kursor zaniżałby ją przy każdym
kolejnym dociąganiu starszych pozycji, a wartość trafia do pola całkowitej liczby, które
warstwa wyższa pokazuje jako rozmiar całej historii, także w ostrzeżeniu przed
nieodwracalnym czyszczeniem okna w całości.

Strona wykazu nie przecina grupy pozycji o jednym znaczniku czasu. Porządek wykazu ma dwa
klucze — czas utworzenia i identyfikator wiersza malejąco — a kursor kontraktu niesie tylko
pierwszy z nich, milisekundy. Gdy granica strony wypada w środku pozycji o tym samym
znaczniku (zwykła tura: pytanie i odpowiedź powstają w tej samej milisekundzie), strona
następna pytana warunkiem ostro mniejszym przeskoczyłaby rodzeństwo bezpowrotnie. Zamiast
dokładać kontraktowi pole rozstrzygające remis, strona domyka się po stronie repozytorium:
po pobraniu limitu wierszy dobierane jest jeszcze całe rodzeństwo ostatniego z nich. Strona
bywa więc odrobinę dłuższa od limitu, za to kursor zawsze pada między grupami.

## budowa/server/internal/dane/biblioteka_odwolania.go
Wykaz odwołań do treści służy sprzątaniu magazynu treści przy starcie rdzenia,
które odczytuje z niego zbiór blobów wciąż powiązanych z bazą. Plik i wersja
trzymają treść przy życiu niezależnie od siebie: blob porzucony przez plik
bieżący może pozostawać treścią wersji historycznej, do której przywrócenie
wersji ma prawo wrócić. Dwa osobne odczyty wymagałyby scalenia po stronie
rdzenia, a pomyłka na scalaniu skasowałaby treść nie do odzyskania — stąd
jedno zapytanie łączące obie tabele przez UNION.

## budowa/server/internal/dane/tlumaczenie_segmentacja.go
Segmenty okna zapisywane są zawsze kompletem: UstawSegmentyOkna wymienia cały wykaz okna w jednej transakcji. Scalenie dwóch segmentów przesuwa numery wszystkich następnych, więc zapis punktowy musiałby i tak dotknąć całego wykazu — tylko w kilku osobnych transakcjach, z oknem, w którym numeracja jest podwójna albo dziurawa.

## budowa/server/cmd/danaco-console/uruchomienie/rozgalezienie.go
Rola hub uruchamia wyłącznie tor interfejsu: nasłuch transportu na porcie
rdzenia, ze strumieniami procesu pozostawionymi nietkniętymi. Rola agent
uruchamia wyłącznie tor wykonawczy: żądania kontraktu przychodzą strumieniem
wejścia, odpowiedzi idą strumieniem wyjścia, żaden port nie jest zajmowany.
Rola all uruchamia oba tory naraz, z torem interfejsu jako wiodącym. Rola
spoza tego katalogu nie zatrzymuje procesu, tylko schodzi na zachowanie roli
all, ponieważ brak rozpoznanego ustawienia ma dawać pracę, nie odmowę.

## budowa/server/internal/dane/tlumaczenie_tresc.go
Ton jest kolumną panelu, nie osobnym bytem: żądanie ustawienia tonu pisze przez metodę tego pliku, bo pole ton mieszka w tabeli panel_tlumaczenia (migracja_053_tlumaczenie.sql). Każda z czterech metod przy nieznanym kodzie panelu wraca ErrBrakWiersza, ponieważ cicha zgoda na zmianę bytu, którego nie ma, byłaby potwierdzeniem czynności, która się nie odbyła — każda metoda sprawdza liczbę wierszy dotkniętych zapisem. Czas jest liczbą milisekund epoki, wzorem reszty modułu Translate. UstawTlumaczenie przyjmuje treść i odwołanie jako wskaźniki: pusta wartość zostawia kolumnę bez zmiany, bo zgłoszenie korekty może nieść samą treść krótką albo samo odwołanie do pliku, zależnie od rozmiaru tekstu.

## budowa/server/cmd/danaco-console/uruchomienie/tor.go
Gdy pole Wejscie pozostaje nieustawione, tor wykonawczy kończy pracę od razu,
nie zrywając przy tym toru interfejsu — proces roli łączącej oba tory ma dalej
obsługiwać interfejs mimo braku podłączonego wejścia.

## budowa/server/internal/dane/tozsamosc.go
Katalog jest sterowany danymi: kilkanaście kategorii to kilkanaście wierszy, nie kilkanaście gałęzi w kodzie, wzorem katalogu akcji. Repozytorium wyłącznie czyta katalog: wiersze wnosi zaczyn migracji 015, a kolejność składania warstw rozstrzyga warstwa wyższa. Treść kategorii mieszka w osobnej tabeli dokument_tozsamosci. Zmiana katalogu jest zmianą danych migracji, nie czynnością kontraktu; treść kategorii zapisuje operator oknem konfiguracji. Porządek zwracanych kategorii jest jednoznaczny, bo składacz promptu ma dawać bajtowo ten sam wynik przy tej samej konfiguracji; kolejność warstw według krytyczności nakłada warstwa wyższa, bo to ona zna silnik nakładki.

## budowa/server/cmd/danaco-narzedzia/stdio/obsluga.go
Zatrzymanie serwera rozpoznaje się między wywołaniami. Drogą właściwą wyjścia
serwera MCP jest zamknięcie strumienia wejścia przez proces modelu, a nie
sygnał — sygnał dochodzi do procesu rodzica, który ten strumień zamyka.
Wiersz nieczytelny dostaje błąd protokołu, metoda nieznana błąd metody,
a odmowa narzędzia wynik oznaczony jako błędny; żaden z tych przypadków
osobno nie kończy pracy serwera.

## budowa/server/internal/dane/tozsamosc_tresc.go
Oś mówi, dla czego treść obowiązuje — dla platformy, dla wskazanego modelu albo dla wskazanego konta. Która oś wygrywa, rozstrzyga warstwa wyższa; repozytorium wyłącznie oddaje wiersze i zapisuje zmianę operatora. Kategoria spoza katalogu jest błędem wskazania, nie awarią, bo katalog jest zamkniętym zbiorem wierszy migracji.
## budowa/server/internal/dane/przegladarka_okno.go
Sprawdzenie znajomości okna operacyjnego jest osobnym pytaniem od odczytu migawki strony, ponieważ wykazy źródeł
i notatek muszą rozróżniać dwie sytuacje wyglądające identycznie przy pustym wyniku: okno przeglądania istnieje,
ale nic w nim jeszcze nie zebrano, oraz okno, którego moduł nigdy nie widział. Bez osobnego sprawdzenia obie
sytuacje kończyłyby się tą samą pustą tablicą i odmową nie do odróżnienia od wyniku pustego.

Śladem obecności okna jest każda z trzech tabel modułu, nie tylko migawka strony. Okno staje się znane pierwszą
nawigacją, ale równie dobrze pierwszym zapisanym źródłem albo notatką — żadna z tych czynności nie wymaga
poprzedniczki. Pytanie ograniczone wyłącznie do tabeli migawek odrzuciłoby jako nieznane okno zasilone samym
źródłem, co byłoby błędem.

Pytanie nie sięga katalogu okien operacyjnych, ponieważ adapter modułu ma jedną zależność — własne repozytorium.
Znajomość okna oznacza więc znajomość modułu przeglądania, a nie istnienie wiersza w tabeli okien.

## budowa/server/internal/dane/biblioteka_indeks_tresci.go
Wiersz pliku, jego filtr i odczyt leżą w library.go — ten plik dokłada tam
jedną klauzulę i jedno polecenie zapisu. Indeks powstaje z treści leżącej
w magazynie na dysku, a nie z kolumny tabeli, i rządzi się własnymi regułami:
co się indeksuje i jak fraza staje się zapytaniem FTS5. Bajty pliku leżą
wyłącznie w magazynie treści rdzenia, w kolumnie tresc_odwolanie; w indeksie
leży wyciąg tekstowy treści bieżącej, obcięty granicą po stronie rdzenia
i służący wyłącznie odnajdywaniu. Brak wiersza indeksu obniża trafność
wyszukiwania i nic poza tym, dlatego zapis indeksu nie wywraca wgrania pliku.

Fraza wyszukiwania jest daną, nie składnią: w przeglądarce biblioteki wpisuje
się słowa, nie wyrażenie FTS5, a znaki cudzysłowu, gwiazdki, minusa, nawiasu
i dwukropka mają w tej składni znaczenie — surowa fraza z nawiasem wywracałaby
wyszukiwanie błędem składni zamiast oddać zero trafień. Cytowanie całej frazy
z podwojeniem cudzysłowu wewnętrznego daje wyszukiwanie sekwencji słów, nie
wyrażenia logicznego. Gwiazdka na końcu trafia przedrostkiem ostatniego słowa
w słowo dokończone, więc trafienia widać przed dokończeniem wpisywanej frazy.

## budowa/server/internal/models/adapter_obrazy_test.go
Sprawdziany mierzą to, czego nie da się zobaczyć po stronie wołającego: kształt żądania, które naprawdę wyszło na sieć. Kanał, który przyjmie materiał operatora i wyśle samo polecenie, oddaje obraz wygenerowany od zera — wygląda to na powodzenie, a jest podmianą materiału bez ani jednego słowa. Dlatego sprawdziany stawiają zaślepkę punktu końcowego i czytają ciało żądania.

## budowa/server/internal/dane/tozsamosc_wyliczenia.go
Kontrakt daje słownik przekładu bazy wyłącznie dla rodzaju konta. Dla warstwy, trybu i osi tożsamości takiego słownika nie ma, bo kolumny trzymają wartości kontraktu wprost — drugie nazewnictwo byłoby drugim źródłem prawdy. Plik nie tłumaczy więc nazw: sprawdza, czy wartość kolumny należy do zbioru kontraktu, i uzupełnia wartość domyślną tam, gdzie wartości nie podano.
## budowa/server/internal/dane/poczta_skrzynki_slad.go
Ślad wysyłki powstaje bezwarunkowo, ponieważ wysłanie listu jest jedyną czynnością tego modułu, której skutek
wychodzi poza system i której nie da się cofnąć. Zapis jest jawny, trwały i niezależny od tego, czy okno rozmowy
jeszcze istnieje — pozwala odtworzyć, co wyszło ze skrzynki i kiedy.

Wysyłka nieudana także jest wierszem tabeli: kolumna powodzenia niesie wartość fałszywą razem z treścią błędu.
Reguła sprawdzająca schematu pilnuje zgodności obu kolumn — wysyłka udana nie może nieść błędu, nieudana musi go
nieść. Dziennik zawierający same powodzenia sugerowałby błędnie, że wszystko doszło do adresata.

Treść listu trafia do śladu w całości, ale bajty załączników nie: leżą w osobnym magazynie pod sumą kontrolną,
a ślad wskazuje je wyłącznie nazwami umieszczonymi w treści.

## budowa/server/cmd/danaco-narzedzia/stdio/protokol.go
Kształt uruchomienia serwera odpowiada wpisowi mcpServers składanemu przez
most MCP rdzenia — typ stdio z poleceniem i argumentami — więc proces modelu
uruchamia ten serwer dokładnie tak, jak uruchamia most konsoli.

Wersja serwera odpowiada wersji produktu. Pakiet nie sięga po stałą wersji
rdzenia, ponieważ import rdzenia wciągnąłby do binarium serwera narzędzi
całą trwałość wraz z bazą.

Kontekst wchodzi parametrem metody Narzedzia, ponieważ w zasięgu eksperta
złożenie wykazu narzędzi pyta rdzeń o jego definicję. Wykaz nie jest więc
czynnością czysto obliczeniową i nie ma prawa przeżyć zatrzymania procesu.

## budowa/server/internal/models/definicja.go
Rejestr kanałów powstaje z wierszy tabeli kanal_modelu w czasie działania aplikacji, więc nowy kanał
znaczy nowy wiersz danych, nie nowy typ w kodzie. Pole PoswiadczenieOdwolanie jest odwołaniem do
danych dostępowych — nazwą zmiennej środowiskowej albo pozycji magazynu; sekret nie jest
przechowywany ani w bazie, ani w repozytorium.

Metoda identyfikatorKontraktu musi zwracać kod kanału, nie numer wiersza. Rejestr indeksuje kanał pod
obydwoma kluczami, więc odczyt działa tak czy inaczej — ale warstwa danych zna wyłącznie kod, więc
numer wiersza dałby klientowi identyfikator, którym nie da się nic zrobić. Numer wiersza jest wartością
zapasową wyłącznie dla wiersza bez kodu.

## budowa/server/internal/dane/design_zasoby.go
Zasoby oddają liczbę wszystkich wyników oddzielnie od stronicowanej listy: zapytanie liczące i zapytanie stronicujące stosują ten sam zestaw warunków WHERE, żeby panel pokazał licznik bez rozjazdu wobec przyciętej listy. Filtr po etykietach jest koniunkcją — zasób musi nieść wszystkie wskazane etykiety naraz, nie choć jedną, dlatego zapytanie wymaga pełnego zestawu przez HAVING COUNT(DISTINCT etykieta) równe liczbie etykiet filtra. Etykiety zasobu są wymianą całego zestawu, nie dokładaniem różnicy: usunięcie zastanych i wstawienie nadesłanych zachodzi w jednej transakcji, żeby zasób nie pozostał przejściowo bez etykiet przy błędzie w trakcie operacji. Usunięcie zasobu idzie po identyfikatorze zewnętrznym; powiązane etykiety znikają kaskadowo przez klucz obcy z ON DELETE CASCADE, a warstwy kompozycji wskazujące ten zasób pozostają nietknięte, ponieważ ich kolumna wskazująca zasób nie niesie klucza obcego.
## budowa/server/internal/dane/poczta_wyzwalacze.go
Odczyt idzie osobnym interfejsem, nie dopiskiem do repozytorium automatyk ogólnego przeznaczenia, z tego samego
powodu co pętla wykonawcza automatyk: interfejs mieszka razem z definicją automatyki, a rdzeń sięga po repozytorium
wyzwalaczy poczty asercją typu na repozytorium automatyk.

Jest to jedna tabela, ale drugie pytanie względem zapisu harmonogramu. Wiersze wyzwalaczy zapisuje wyłącznie warstwa
harmonogramu automatyk; tutaj leży odczyt odwrotny — od skrzynki do automatyki — którego moduł planujący nigdy
nie zadaje. Wyrażenie gwiazdki i wyrażenie puste znaczą obie „każda skrzynka", a droga kontraktu zapisuje gwiazdkę,
ponieważ ustawienie harmonogramu pomija wyzwalacze z wyrażeniem pustym.

## budowa/server/internal/dane/urzadzenia.go
Urządzenie to maszyna z klientem albo z katalogiem udostępnionym modelowi. Punkt dostępu rodzaju localDirectory bez wskazania urządzenia nie przechodzi więzu schematu, bo ścieżka lokalna ma znaczenie tylko na jednej maszynie. Pole biezace nadaje wyłącznie ZapewnijBiezace, bo tylko ono potrafi najpierw zdjąć oznaczenie z pozostałych wierszy — baza dopuszcza jedno takie urządzenie.

## budowa/server/cmd/danaco-narzedzia/stdio/metody.go
Kształt pozycji wykazu narzędzi składa funkcja pomocnicza narzędzi, nie
wykazNarzedzi: licznik wykazu mierzy wagę tej samej odpowiedzi w bajtach,
a dwa miejsca składające ją osobno dałyby pomiar czegoś innego, niż dostaje
model. Funkcja wykazNarzedzi zostaje więc wyłącznie zamówieniem wykazu.

Model ma przeczytać odmowę narzędzia i błąd rdzenia oraz poprawić wywołanie,
zamiast dostać usterkę połączenia — dlatego wracają wynikiem oznaczonym jako
błędny. Błędem protokołu zostaje wyłącznie treść żądania, której nie da się
odczytać.

## budowa/server/internal/models/fragment.go
Kontrakt nie opisuje kształtu ładunku obrazu, ponieważ pole `data` fragmentu rodzaju `image` jest
surowym JSON-em — kształt TrescObrazu mieszka więc w tym pakiecie. Nazwy pól są przepisane z miejsc,
gdzie kontrakt już opisuje treść binarną (`mimeType`, `contentBase64`, `uri` — LibraryFileUpload,
DesignAsset), żeby odbiorca składający zasób wizualny nie tłumaczył nazw po drodze. Dostawcy zgodni
z OpenAI Images oddają albo bajty w Base64, albo odsyłacz w Adres, zależnie od `response_format` —
fragment powtarza to, co przyszło, i niczego nie dosypuje.

## budowa/server/internal/dane/design_zasoby_zmiany.go
Obie metody wykonują jedno polecenie SQL bez transakcji, w odróżnieniu od zmiany etykiet, która składa się z dwóch poleceń i wymaga transakcji, żeby stan pośredni nie był widoczny. Usunięcie zasobu nie rusza bajtów w magazynie: blob leży pod sumą swojej zawartości, więc dwa zasoby o tej samej treści mogą dzielić jeden plik, a skasowanie pliku przy usunięciu jednego z nich odebrałoby treść drugiemu; o bajtach rozstrzyga warstwa rdzenia, tu znika sam wiersz. Ustawienie ulubionego traktuje brak wiersza do zmiany jako wynik neutralny, nie błąd, bo wołający ma klucz wiersza już po odróżnieniu zasobu nieznanego. Usunięcie zasobu opiera zwracany wynik na liczbie zmienionych wierszy, nie na powodzeniu polecenia, ponieważ usunięcie zera wierszy jest dla bazy sukcesem — powtórne usunięcie tego samego identyfikatora zwraca wynik ujemny, nie błąd.

## budowa/server/internal/dane/automations_kolejka.go
Plik nie jest drugim silnikiem kolejek: cykl życia zlecenia — stany, werdykt,
bieg naprawczy — prowadzi wyłącznie silnik automatyk nad repozytorium kolejek.
Tutaj leżą dwie operacje ułożenia pozycji, która jest wcześniej i w której
kolejce stoi; żadna z nich nie zmienia stanu pozycji ani nie posuwa jej
naprzód. Zapis idzie przez ten sam dziennik akcji kolejki co reszta ruchu
pozycji, więc ślad zostaje w jednym miejscu.

## budowa/server/internal/dane/design_zetony.go
Zapis zestawu żetonów jest zawsze pełny, ponieważ kontrakt zapisu nie zna trybu częściowej zmiany, a panel projektowy nadsyła stan docelowy całego drzewa. Zapis usuwa zastane żetony i wstawia nadesłane w jednej transakcji tym samym wzorcem co przy etykietach zasobu — inaczej żeton skasowany w panelu zostawałby w bazie i wracał przy następnym odczycie. Klucz zestawu odczytuje się przed podmianą żetonów, bo zestaw mógł dopiero powstać w tej samej transakcji, a tabela żetonów wymaga gotowego klucza wiersza.

## budowa/server/cmd/danaco-console/main.go
Cały odczyt zmiennych środowiska prowadzi pakiet konfiguracja, którego wykaz
nazw pokrywa się ze wzorcem pliku przykładowego środowiska — żadna zmienna
środowiska nie jest czytana po nazwie w punkcie wejścia. Poza kompozycją nie
ma tu logiki, typów ani obsługiwaczy komend: to, co robi rdzeń, mieszka
w warstwie rdzenia, a wybór torów w warstwie uruchomienia.

Tryb wypisania wykazu zależności zewnętrznych stoi przed odczytem
konfiguracji, ponieważ zestaw flag konfiguracji odrzuciłby ten argument jako
nierozpoznany. Wykaz idzie na wyjście standardowe, skąd konsumuje go
prowizjonowanie serwera, żeby nazwy pakietów miały jedno źródło.

Kontrola integralności pliku bazy i kontrola kluczy obcych wykrywają osobno
uszkodzenie pliku i wiersze osierocone, oba tuż po otwarciu i migracjach.

Wpięcie toru do hosta zdalnego czyta tabelę hosta zdalnego i ustawienie hosta
wykonania przez uchwyt podany w punkcie wejścia; sam nie otwiera bazy, bo
właścicielem jedynej puli połączeń procesu jest kompozycja. Bez tego wpięcia
każda droga toru odmawia z powodu jego braku; z nim odmowy zostają tylko tam,
gdzie brakuje wskazania hosta albo zgody na maszynę.

## budowa/server/internal/session/obieg.go
Bieg naprawczy nie ma limitu obiegów. Zamiast bramy licznikowej wchodzi przejrzystość:
licznik obiegów, wykrywanie braku postępu i jawny, nazwany warunek zatrzymania.
Zatrzymanie nigdy nie jest ciche, zawsze niesie rozpoznany powód przekazywany
obserwatorom pętli. Powody dzielą się na dwie klasy rozstrzygane maszynowo, nie
opisowo: trzy zastają pracę przerwaną i bieg podejmuje z nich wyłącznie operator
poleceniem wznowienia; czwarty, ukończenie z wynikiem, zastaje pracę zrobioną
i bieg podejmuje się z niego sam, gdy praca dostanie ciąg dalszy. Dzięki tej
różnicy układ złożony przez operatora poznaje wynik zadania bez pytania.

## budowa/server/internal/dane/urzadzenia.go (uzupełnienie)
ZapewnijBiezace zakłada albo odświeża wiersz maszyny, na której działa rdzeń, i przenosi na nią oznaczenie maszyny bieżącej. Wywołanie powtórzone tymi samymi znamionami nie tworzy drugiego wiersza.
## budowa/server/internal/dane/terminal.go
Repozytorium nie prowadzi procesów: uchwyt do procesu działającego i jego drzewa potomstwa ma wyłącznie rdzeń,
ponieważ tylko on potrafi proces zakończyć. W repozytorium zapisuje się to, co po procesie zostaje: polecenie,
inicjator, kod wyjścia i czasy. Dzięki temu podgląd procesów pokazuje także procesy zakończone, których żywy
stan rdzenia już nie przechowuje.

Pole CelZdalny niesie adres powłoki zdalnej w postaci użytkownik i host albo alias konfiguracji zdalnego dostępu
maszyny rdzenia. Nie jest poświadczeniem — jest tą samą wartością, co widnieje w książce hostów — dlatego ma
własną kolumnę, inaczej niż zmienne środowiska karty wprowadzone późniejszą migracją schematu.

Procesy osierocone to procesy pozostawione w stanie działania przez poprzedni bieg rdzenia. Po ponownym
uruchomieniu rdzeń nie ma już do nich uchwytu, więc wykazywanie ich jako czynnych byłoby niezgodne ze stanem
faktycznym.

## budowa/server/cmd/danaco-console/uruchomienie/tor_wykonawczy.go
Tą drogą pracuje agent lokalny: wykonuje pracę na urządzeniu użytkownika, nie
otwierając gniazda ani nie zajmując portu. Żądanie nieczytelne nie kończy
toru — rdzeń odsyła na nie odpowiedź z kodem błędu kontraktu.

## budowa/server/internal/dane/urzadzenia_zapis.go
Oznaczenie maszyny bieżącej ma jednego pisarza: ZapewnijBiezace. Unikatowy indeks częściowy idx_urzadzenie_biezace dopuszcza najwyżej jedno urządzenie z biezace ustawionym, więc nadanie oznaczenia wymaga wcześniejszego zdjęcia go z pozostałych wierszy. Dodaj i Aktualizuj kolumny biezace nie ruszają. Powtórzone rozpoznanie tej samej maszyny ma zaktualizować wiersz, a nie założyć drugi — kolumny nazwa i zaufane zostają nietknięte, bo pierwsza bywa zmieniona przez operatora, a druga bywa świadomie odebrana.

## budowa/server/internal/dane/developer.go
Repozytorium nie dotyka plików na dysku: treść pliku roboczego żyje w katalogu roboczym okna i czyta ją rdzeń, a warstwa danych zapisuje wyłącznie migawkę zakładaną na wyraźne żądanie zapisu wersji, czyli stan, którego po nadpisaniu pliku na dysku już nie ma. Repozytorium nie prowadzi też samego budowania — uchwyt do biegnącego procesu ma rdzeń, tutaj zostaje wyłącznie ślad: zadanie, kod wyjścia, czasy i ogon dziennika. Przebieg już domknięty zostaje przy próbie ponownego domknięcia bez zmiany, żeby pierwszy prawdziwy kod wyjścia nie został nadpisany przez późniejsze przerwanie. Przebiegi zostawione w stanie running po poprzednim uruchomieniu rdzenia przechodzą po restarcie w stopped, ponieważ rdzeń nie ma już do nich uchwytu i wykazywanie ich jako czynnych byłoby nieprawdą.

## budowa/server/internal/session/polecenie.go
Pakiet session nie wie, jak zbudować wiersz poleceń kanału modelu, ponieważ to
należy do warstwy kanału. Rozdzielenie idzie przez interfejs, nie przez wywołanie
w głąb cudzego pakietu. Polecenie jest wspólnym kształtem, nie punktem uruchomienia:
wypełnia je rdzeń poprzez moduły Terminal i Developer, a wykonuje warstwa kanału
przez port Uruchamiacz. Sesja nie buduje polecenia i sama go nie wykonuje.
## budowa/server/internal/dane/schedule.go
Zapis harmonogramów, odczyt po automatyce, wyzwalacze i budzik stoją w osobnym pliku warstwy danych automatyk
i nie powtarzają się tutaj. Drogi odczytu są dwie, bo pytający bywa różny: ustawienie harmonogramu zna automatykę
i pyta o jej harmonogram kolumną klucza obcego automatyki, natomiast odczyt harmonogramu zna albo identyfikator
zewnętrzny samego harmonogramu, albo nie zna żadnego wskazania i pyta o komplet. Żadnej z tych dróg nie da się
przejechać jednym zapytaniem po kluczu obcym automatyki, więc dochodzą tu dwa osobne zapytania.

Schemat niesie oba potrzebne więzy jednoznaczności: jeden na identyfikator zewnętrzny, dzięki czemu odczyt po
kodzie oddaje co najwyżej jeden wiersz, i drugi na klucz obcy automatyki, dzięki czemu automatyka ma najwyżej
jeden harmonogram.

## budowa/server/internal/models/kanal.go
Rdzeń nie zna żadnego adaptera kanału z osobna — sięga po adaptery wyłącznie przez rejestr. Adapter
kanału głównego (Claude Code CLI) realizuje interfejs Kanal w pakiecie internal/injection; tam mieszka
budowa argumentów wywołania, pula kont i parser strumienia odpowiedzi.

Metoda Wyslij interfejsu Kanal nadaje jako pierwszy fragment strumienia prowenancję wywołania. Zwrócony
błąd dotyczy wyłącznie tego wywołania; adapter błędu nie zamienia sam na fragment — robi to
Rejestr.Wyslij, dzięki czemu strumień nie niesie dwóch fragmentów błędu o tej samej przyczynie.

Stała AdapterObrazy stoi obok AdapterAPI, a nie zamiast niego: rodzaj kanału (`api`) mówi, jak kanał
rozmawia, a adapter mówi, co oddaje.

## budowa/server/internal/dane/ustawienia_osi.go
Oś jest prostopadła do poziomu: poziom mówi, jak wąsko obowiązuje wartość, oś mówi, dla czego — dla platformy, dla modelu albo dla konta. Wersje metod bez osi opisują oś platformy i wywołują dokładnie ten sam kod, więc rozstrzyganie ma jedną implementację, nie dwie. Kolejność osi i ich pierwszeństwo trzyma pakiet konfiguracji rdzenia. Tabela os_zasiegu jest wyłącznie więzem klucza obcego dla kolumny ustawienie.os.

## budowa/server/internal/dane/developer_odczyt.go
Limit wykazu wchodzi zapytaniem jako parametr przygotowanego polecenia, nie sklejaniem tekstu SQL. Wartość niedodatnia oznacza wykaz pełny dzięki wyrażeniu warunkowemu w zapytaniu, więc jedno przygotowane polecenie obsługuje oba przypadki bez rozgałęzienia w kodzie, a liczba z zewnątrz nigdy nie trafia do treści zapytania.

## budowa/server/internal/session/wstrzymanie_windows.go
Windows nie ma dla obcego procesu odpowiednika sygnałów uniksowych wstrzymania
i wznowienia. Wstrzymanie idzie tam per wątek, a uchwyt, którym rdzeń obejmuje
całe drzewo procesów, takiej czynności nie zna. Przejście po wątkach wszystkich
procesów drzewa nie jest odpowiednikiem, bo wątek utworzony w trakcie czynności
zostałby biegnący, więc wstrzymanie znaczyłoby coś innego niż na systemach
uniksowych. Dlatego czynność zgłasza się jako niewspierana, zamiast robić coś
podobnego: kontrakt komendy zgłoszenia wstrzymania ma na to osobne pole,
a proces zostaje nietknięty, więc operator dostaje prawdę o możliwościach
maszyny zamiast odpowiedzi udanej, po której proces dalej zajmuje procesor.

## budowa/server/internal/models/konto.go
Kanał nadaje fragment MetadaneKonta także w chwili przełączenia konta w trakcie sesji — rotacja konta
ma być widoczna wołającemu, nie milcząca.
## budowa/server/internal/dane/rozmowa_bloki.go
Typ wiadomości niesie pole pojemne metadanych w postaci surowego JSON; bloki wracają w nim pod własnym kluczem,
obok atrybucji i zużycia, które warstwa danych rozmowy składa z osobnych kolumn. Klient odtwarza z nich tok
rozumowania, narzędzia i pochodzenie wpisu tą samą logiką, którą składa turę żywą po stronie interfejsu.

Kształt bloku w metadanych powtarza kształt fragmentu strumienia zdarzeń po nazwach pól kontraktu, z dodaną
chwilą powstania. Dzięki temu klient nie potrzebuje drugiego słownika rodzajów bloku ani drugiego czytnika
ładunku — czyta bloki zapisane i bloki strumieniowe tym samym kodem.

## budowa/server/internal/dane/developer_warsztat.go
Byty tego pliku nie leżą w pamięci rdzenia, bo sesja debugowania, połączenie z bazą i biegnący skan są stanem żywym z uchwytami do procesów, gasnącym razem z nimi — taki stan w bazie nie ma miejsca. W bazie leży to, co osoba pracująca w edytorze ułożyła i czego nie może stracić przy zamknięciu okna: postawiony punkt przerwania, zapisane zapytanie, opisane połączenie, wynik pomiaru.

## budowa/server/cmd/danaco-narzedzia/main.go
Poza kompozycją nie ma tu wykazów, komend ani protokołu — wszystko to mieszka
w warstwie narzędzi rdzenia i w podpakiecie stdio. Dziennik idzie na wyjście
diagnostyczne, bo wyjście standardowe należy w całości do protokołu MCP —
jeden obcy wiersz na standardowym wyjściu zerwałby rozmowę z procesem modelu.

Przełącznik zasięgu czyta się raz, przy uruchomieniu, z wpisu okna ułożonego
przez rdzeń — proces modelu startuje z gotowym wykazem i nie ma czym go
poszerzyć w trakcie rozmowy. Kod eksperta rdzeń bierze z pola ustawień agenta
okna. Przełącznik dołożeń sesji składa strona rdzenia, stamtąd bierze się
jego nazwa; analiza argumentów kończy proces na przełączniku nieznanym, więc
odczyt dołożeń musi stać, zanim tamta strona zacznie argument dokładać.

Brak okna nie zatrzymuje serwera: model podaje okno sam albo rdzeń odmawia.
Cudzego okna serwer nie podstawi nigdy, więc milczenie jest tu bezpieczne;
głośne pozostaje w dzienniku, bo wpis okna zawsze je niesie.

Liczby zestawu przy starcie dotyczą zestawu przed zawężeniem, ponieważ
gniazdo do rdzenia jest leniwe i definicji eksperta nie da się mieć w tej
chwili — cenę zestawu eksperta melduje dobór przy pierwszym wykazie narzędzi.

Cichy odrzut dołożenia bez zawężenia zostawiłby narzędzie na wykazie sesji
bez śladu, że w tej turze nic nie zmieniło, dlatego dziennik go odnotowuje.

Rozjazd dwóch przełączników zasięgu jest możliwy w obie strony. Kod eksperta
bez zasięgu eksperta podnosi zasięg do zasięgu eksperta, z podniesieniem
odnotowanym w dzienniku. Zasięg eksperta bez kodu eksperta nie ma o kogo
zapytać rdzenia i schodzi do zasięgu okna, ponieważ zawężenie bez eksperta
nie jest zawężeniem, tylko obietnicą bez pokrycia.

## budowa/server/internal/transport/rozgloszenie.go
Synchronizacja wielourządzeniowa nie ma własnego protokołu: nośnikiem jest
zdarzenie właściwe zmienionemu obszarowi, rozgłoszone drogą rozgłoszenia kopert.

## budowa/server/internal/dane/asystent_polecenia.go
Polecenie przyjęte bez śladu w dzienniku, gdzie zlecenie jest, a rozmowa nie
ma pierwszej linii, albo ślad bez zlecenia, gdzie wpis wisi na kodzie, którego
zlecenie nigdy nie powstało, jest stanem połowicznym — stąd zapis w jednej
transakcji, tym samym wzorem co przy zapisie kompozycji i zapisie kroków
gdzie indziej w module. Ten plik nie rozpoznaje mowy: żądanie niesie odnośnik
do nagranego już pliku albo tekst poprawiony przez operatora, a zapis idzie
dosłownie do odpowiednich kolumn. Brak obu jest błędem żądania, nie pustym
zapisem, ponieważ kolumna treści jest wymagana — wpis dziennika bez treści
nie opisuje niczego, co się wydarzyło.

## budowa/server/internal/dane/wiadomosci.go
Pole IdentyfikatorZewnetrzny wiąże wiersz z wiadomością rdzenia, którą klient zna pod identyfikatorem tekstowym; po nim odnajduje się wiersz odpowiedzi modelu, gdy strumień się domyka. Pole Zalaczniki niesie wykaz odwołań do załączników w postaci tablicy JSON napisów, dokładnie tak, jak brzmi pole attachments kontraktu. Wartość nil znaczy pustą kolumnę, czyli że baza nic o załącznikach tej wiadomości nie wie; to nie to samo co pusta tablica, która znaczy sprawdzone: nie było żadnych.

## budowa/server/internal/session/uruchamiacz.go
Uruchamianie procesu ma w drzewie jedno miejsce, warstwę kanału, bo tylko ona
zna wiersz poleceń, rotację kont i kształt strumienia. Pakiet session nie
buduje ani nie startuje procesu, także procesu okna komunikacji: bierze uchwyt
procesu gotowego i dokłada to, co należy do sesji. Interfejs zostaje w tym
pakiecie, choć sam session po niego nie sięga, ponieważ definiuje go strona
znająca kształt uruchomienia okna, a wypełnia warstwa kanału; sięgają po niego
moduły Terminal i Developer w rdzeniu, bo one uruchamiają procesy okna i one
obejmują je drzewem.

## budowa/server/internal/dane/developer_warsztat_odczyt.go
Zawężenia wykazów wchodzą do zapytań jako parametr przygotowanego polecenia, nie sklejaniem tekstu SQL. Wartość pusta oznacza brak zawężenia, więc jedno przygotowane zapytanie obsługuje zarówno wykaz pełny, jak i zawężony, a wartość z zewnątrz nigdy nie trafia bezpośrednio do treści polecenia — ten sam wzorzec obsługuje limit wykazu w odczycie warsztatu.

## budowa/server/internal/session/wstrzymanie.go
Rdzeń umiał proces wyłącznie zakończyć, więc zadanie długie, na przykład
budowanie, wielka migracja albo transkodowanie, dało się tylko ubić, czyli
wyrzucić wykonaną pracę. Wstrzymanie oddaje maszynę bez utraty postępu: proces
przestaje dostawać czas procesora, jego pamięć i otwarte pliki zostają,
a wznowienie podejmuje pracę w miejscu, w którym stanęła. Czynność ta nie jest
dostępna na wszystkich systemach jednakowo: na systemach uniksowych są do niej
sygnały wysyłane całej grupie procesów, dokładnie tej samej, którą obejmuje
ubicie, a Windows nie ma dla obcego procesu odpowiednika, bo wstrzymanie idzie
tam per wątek, czego uchwyt obejmujący całe drzewo nie zna. Kontrakt komendy
wstrzymania ma dlatego osobne pole wsparcia platformy, a nie tylko odpowiedź
udaną albo nieudaną, bo brak wsparcia systemu to co innego niż niepowodzenie.
Rozstrzygnięcie stoi w plikach zależnych od platformy, tak samo jak ubicie.
## budowa/server/internal/dane/przejecie_sterowania.go
Ślad przejęcia sterowania leży w istniejącym dzienniku akcji okna, bez osobnej tabeli, ponieważ przejęcie
sterowania jest akcją wykonaną na oknie, a te mają już jeden wspólny dziennik. Kolumna nazywająca akcję jest
w tej tabeli wolnym tekstem, bez ograniczenia sprawdzającego i bez klucza obcego do katalogu akcji, więc oba
rodzaje wpisu mieszczą się w niej bez zmiany schematu — druga tabela na to samo zdarzenie byłaby drugą prawdą
o jednym fakcie.

Plik dokłada wyłącznie odczyt zawężony do tych dwóch rodzajów wpisu, zapis idzie istniejącą funkcją zapisu akcji.
Odczyt musi być osobny, ponieważ ogólny odczyt dziennika akcji okna oddaje dziennik w całości i z limitem: okno
z dużą liczbą zwykłych akcji zepchnęłoby przejęcia poza limit, czyli historia sterowania milczałaby dokładnie
tam, gdzie jest najbardziej potrzebna.

Parametry i wynik zostają surowe, tak samo jak w sąsiednim pliku dziennika akcji okna: kształt obu kolumn zależy
od rodzaju akcji, którego warstwa danych nie zna. Ten plik ich nie rozbiera — oddaje wpis w tej samej postaci,
co dziennik akcji, a rozbiór wykonuje warstwa wyższa.

Wąski kontrakt odczytu śladu sterowania pozwala warstwie wyższej sięgnąć po tę jedną metodę asercją typu, bez
dopisywania linii do szerokiego kontraktu obszaru okien.

## budowa/server/internal/dane/asystent_profil.go
Profil niesie warstwę promptu zlecenia — zdanie, które mówi modelowi, że jest
klawiaturą operatora, a nie autorem odpowiedzi. To ono rozstrzyga, czy model
sięgnie po narzędzia platformy, czy odpisze tekstem. Reszta kolumn to nastawy
tury, które w innym razie bierze wiersz okna. Plik ma sam odczyt, bez zapisu:
zakładanie profilu, wykaz profili i wskazanie domyślnego to trzy czynności
operatora, a kontrakt nie ma dla nich ani jednej komendy; metoda zapisu bez
wołającego byłaby drogą, której nikt nie przechodzi.

Pola środowiska wykonania i trybu uprawnień niosą wartości kontraktu wprost,
ponieważ kolumna ma na nich warunek CHECK, a rdzeń wkłada je do zapytania
kanału bez przekładu. Pakiet danych nie zależy od pakietu kontraktu, więc
typem jest tu napis; jedynym miejscem, w którym te napisy stają się typami
kontraktu, jest adapter modułu.

Wskaźniki przy warstwie promptu, kanale modelu i głosie syntezy są rozmyślne:
wartość pusta bazy znaczy, że profil nie ma w tej sprawie zdania i wtedy
obowiązuje nastawa okna. Pusty napis znaczyłby, że profil kasuje nastawę
okna — a profil ma zasięg pracy poszerzać, nie zabierać.

Domyślny profil jest co najwyżej jeden, czego pilnuje indeks częściowy bazy.
Ograniczenie liczby wyników do jednego stoi w zapytaniu mimo to, żeby odczyt
nie zależał od tego, czy ten indeks przetrwał każdą przyszłą zmianę schematu.

## budowa/server/internal/dane/workspace.go
Projekt zakłada się przy pierwszym wejściu do niego: kontrakt nie ma komendy zakładającej projekt osobno, a widok pulpitu projektu ma odpowiedzieć, nie odmówić — brak wiersza jest więc stanem początkowym, nie błędem.

## budowa/server/internal/dane/developer_warsztat_zapis.go
Zapisy zbiorcze znalezisk skanu, wyników testów i pokrycia idą jedną transakcją i zaczynają się od usunięcia poprzedniego pomiaru, ponieważ pomiar jest stanem z jednej chwili, nie przyrostem: dopisanie drugiego przebiegu do pierwszego dałoby wykaz, w którym ten sam test stoi dwa razy z dwoma różnymi wynikami, bez sposobu rozstrzygnięcia, który wynik jest aktualny.

## budowa/server/internal/zdalne/nadajnik_powiadomien.go
Pakiet opisuje kształt nadajnika, którego potrzebuje, i nie bierze interfejsu
z pakietu transportu ani z rdzenia, podobnie jak rdzeń opisuje u siebie własny
port nadajnika zamiast importować pakiet transportu. Port nie niesie koperty:
nośnik rozgłoszenia przyjmuje kopertę z typem komunikatu z kontraktu, a kontrakt
nie ma dziś ani jednego zdarzenia powiadomienia, więc zbudowanie koperty z typem
spoza kontraktu wysłałoby komunikat, którego klient nie zna. Dlatego port bierze
sam fakt do przekazania, a złożenie koperty należy do strony znającej kontrakt,
czyli adaptera w rdzeniu wpiętego po wniesieniu zdarzenia powiadomienia. Ten
pakiet prowadzi kolejkę i rozstrzyga, co i kiedy ma polecieć, a nie rozstrzyga,
jakim napisem się to nazywa na łączu. Ponowienie doręczenia kosztuje jeden takt,
natomiast nieprawdziwy zapis stanu doręczenia w bazie kosztuje zaufanie do kanału.
## budowa/server/internal/dane/szukanie_rozmow.go
Szukanie idzie przez dopasowanie pełnotekstowe, wycinek i porządek trafienia, nie przez porównanie tekstu.
Tokenizator indeksu sprowadza większość znaków diakrytycznych do liter podstawowych, ale znak „ł" pozostaje
osobną literą, więc zapytanie bez tego znaku nie trafi w słowo, które go zawiera.

Fraza wyszukiwania trafia do zapytania w cudzysłowie: zapytanie jest frazą dosłowną, a nie wyrażeniem składni
wyszukiwania pełnotekstowego. Ujęcie w cudzysłów, z podwojeniem cudzysłowów wewnętrznych, zdejmuje z frazy
operatory składni, żeby fraza z myślnikiem albo gwiazdką nie wywracała zapytania błędem składni.

## budowa/server/internal/session/wstrzymanie_unix.go
Wstrzymanie samego korzenia zostawiłoby biegnące potomstwo, czyli tę część
pracy, która zwykle zajmuje maszynę, dlatego sygnał idzie do ujemnego
identyfikatora grupy, tak samo jak przy ubiciu procesu. Straż pilnująca
identyfikatora większego niż jeden pilnuje tego samego, co przy ubiciu: sygnał
do minus jeden byłby rozgłoszeniem do wszystkich procesów systemu, a do jeden,
sygnałem do procesu init. Brak procesu o danym identyfikatorze znaczy, że nie
ma już czego wstrzymywać, i nie jest błędem.

## budowa/server/internal/dane/diagnostics.go
Repozytorium nie wytwarza faktów diagnostycznych: nie liczy stanu systemu, nie ocenia błędów i nie wymyśla rekomendacji, wyłącznie przenosi to, co rdzeń mu przekazał. Wzorzec dopasowania w filtrze dziennika dotyczy treści wpisu, a samo dopasowanie regularne wykonuje rdzeń, ponieważ SQLite bez rozszerzenia nie zna operatora dopasowania wyrażeń regularnych. Zapis błędu diagnostycznego rozstrzyga po odcisku: wystąpienie o odcisku już znanym podnosi licznik i przesuwa chwilę ostatniego wystąpienia, a wystąpienie o odcisku nowym zakłada nowy wiersz od zera.

## budowa/server/internal/models/prowenancja.go
Prowenancja jest przejrzystością wywołania, nie bramką — skrót nakładki służy wyłącznie diagnostyce
i nigdy nie dopuszcza ani nie blokuje wywołania. Nazwy pól JSON szkieletu są wspólne
z injection.Prowenancja kanału głównego; klient rozmowy parsuje wyłącznie ten jeden kształt, więc
kanały sieciowe i echo emitują prowenancję o tych samych nazwach pól co kanał główny. Pola specyficzne
kanałów bez procesu (kanał, adapter, adres, ustawienia wywołania, środowisko) jadą jako dodatkowe, poza
wspólnym szkieletem, i nie kolidują z jego nazwami.

Wywołujący pakietu ProwenancjaZapytania dokłada jedynie to, co zna sam adapter: argv albo adres.

Pole PlikUstawien niesie nazwę pola JSON „settings", taką samą jak w kanale głównym. Pole Ustawienia
nosi osobną nazwę pola JSON „callSettings", by nie kolidować ze szkieletowym „settings".

## budowa/server/internal/transport/wykonanie.go
Funkcja stosuje dwa zachowania odporne na awarię: gdy rdzeń nie jest podłączony,
żądanie dostaje zdarzenie nieznanej komendy z kontraktu, a połączenie żyje
dalej; gdy obsługa się załamie, wraca odpowiedź z kodem błędu wewnętrznego
zamiast przerwania procesu, a kolejne żądania są przyjmowane. Straż bramki nie
jest odporna w ten sposób i ma pierwszeństwo przed obydwoma zachowaniami: poza
pętlą zwrotną żądanie z gniazda, które nie przedstawiło tokenu, nie dochodzi do
rdzenia i wraca odmowa braku uwierzytelnienia; na pętli zwrotnej straż jest
wyłączona. Odpowiedź na komunikat niepoprawny strukturalnie różni się od
odpowiedzi na nieznaną komendę kodem błędu walidacji.

## budowa/server/internal/dane/diagnostics_analiza.go
Analiza i jej rekomendacje zapisują się razem albo wcale. Rekomendacja bez
analizy nie ma faktu, z którego wynika, a analiza z połową rekomendacji
kłamie o tym, co z niej wypadło. Jedna transakcja zamyka obie możliwości.

## budowa/server/internal/session/okno.go

Normalizacja w funkcji uzupelnijUstawienia obejmuje wszystkie trzy wyliczenia okna, nie tylko rolę.
Wartość spoza słownika kontraktu, na przykład tryb uprawnień nierozpoznany przez rejestr, trafiłaby
do rejestru pamięciowego bez przeszkód, ponieważ rejestr taką wartość znosi. Wiersz okna komunikacji
w bazie już nie: przekład na kolumnę wyliczeniową odmawia zapisu, bo wartość nie należy do słownika
kontraktu, zapis okna pada, a dziennik rozmowy schodzi całym oknem na bufor pamięci. W takim stanie
odczyt listy wiadomości pokazuje rozmowę z bufora, natomiast wczytanie i usunięcie historii widzą
pustkę w bazie, ponieważ okno bez wiersza nie ma historii ani danych do przycięcia retencją. Wartość
nieznaną normalizacja doprowadza więc do wartości domyślnej, tak samo jak czyni to normalizacja roli
okna: bez odmowy zapisu, za to z oknem, które da się utrwalić. Słowniki wartości pochodzą z kontraktu,
nie z literałów wpisanych lokalnie w kodzie.
## budowa/server/internal/dane/panele.go
Zapis układu sekcji panelu jest całościowy, nie różnicowy: podane sekcje wyznaczają układ panelu, a czego w żądaniu nie ma, tego po zapisie nie ma w bazie. Wynika to z kontraktu — polecenie zapisu niesie samo pole listy sekcji, bez znacznika czynności, więc jedynym czytelnym znaczeniem listy jest układ docelowy. Zapis różnicowy wymagałby, żeby klient wiedział, co w bazie leży teraz; wtedy dwa okna przestawiające ten sam panel rozjechałyby układ, bo każde dopisywałoby swoje do cudzego stanu. Skasowanie starego układu i wpisanie nowego idzie jedną transakcją, bo przerwane w połowie zostawiłyby panel bez sekcji, nieodróżnialny od panelu nigdy nieustawianego. Kolejność nadaje ten plik, nie wołający: baza pilnuje wyłącznie dolnej granicy numeru, a numery 1..N nanosi zapis sekcji, żeby wykaz czytany po numerze porządkowym był tym samym, co wykaz czytany po miejscu na liście.

Odczyt po zapisie idzie po zamknięciu transakcji i po opadnięciu zapisów zbiegłych w czasie: układ obowiązujący to stan, który zastanie następny czytelnik. Dwa okna przestawiające ten sam panel dostają dzięki temu tę samą treść, więc odpowiedź niepodobna do żądania znaczy, że cudzy zapis wszedł po naszym. Licznik zapisów w toku pozwala odczytać układ dopiero wtedy, gdy zapisy zbiegłe w czasie opadły; czekanie ma kres, bo panel przestawiany bez ustanku nie może wstrzymać odpowiedzi w nieskończoność, a odczyt po upływie kresu jest nadal odczytem z nośnika.

## budowa/server/internal/models/rejestr.go
Stan „czynny” zwracany przez Kontrakt nie pokrywa się z wartością aktywny wiersza
źródłowego. Wiersz może mieć aktywny równe jeden, a mimo to nie mieć zbudowanego
adaptera, bo jego rodzaj nie ma fabryki albo fabryka odmówiła budowy — taki wiersz
trafia do Pominiete. Kanał bez adaptera nie jest gotowy do pracy, więc przy
ograniczeniu do kanałów czynnych nie jest zwracany jako czynny; inaczej okno
wskazywałoby kanał, który przy pierwszej próbie odmawia działania.

## budowa/server/internal/dane/diagnostics_bledy.go
Wystąpienie podnosi licznik, nie zakłada wiersza. Ten sam błąd powtórzony
tysiąc razy jest jednym wierszem o tysiącu wystąpień. Gdyby każde wystąpienie
zakładało wiersz, panel błędów pokazywałby ostatnią minutę pracy i gubił błąd
rzadki, a to właśnie rzadki błąd bywa przyczyną awarii.

Stan, priorytet i notatka należą do operatora. Powtórne wystąpienie nie
przestawia ich z powrotem na wartość początkową: rozstrzygnięcie operatora nie
ma prawa zniknąć dlatego, że błąd wystąpił jeszcze raz.

Zliczenie bez ograniczenia liczby wyników istnieje po to, żeby wykaz błędów
mógł podać liczbę całkowitą niezależną od długości zwróconej listy — inaczej
liczba całkowita zawsze równałaby się liczbie oddanych wierszy i panel błędów
nigdy nie dowiedziałby się, że wykaz ucięto na granicy pięciuset pozycji.

## budowa/server/internal/models/rodzaj_fragmentu.go
Rodzaje diagnostyczne prowenancji i konta są w kontrakcie wartościami przelotowymi:
nie mają odpowiednika w kolumnie wiadomosc.rodzaj_tresci, więc nie wchodzą do
WartosciBazyChunkKind i nie trafiają do bazy, której ograniczenie CHECK dopuszcza
wyłącznie rodzaje treści rozmowy. Rodzaje treści rozmowy nie mają tu drugiej
nazwy: bierze się je wprost ze współdzielonego kontraktu, a odwzorowanie na
kolumnę bazy wykonuje pakiet dane.

## budowa/server/internal/models/sciezka_tresci.go
Odczyt wartości ścieżką z dokumentu JSON pozwala adapterowi sieciowemu nie znać
kształtu odpowiedzi żadnego dostawcy: kształt jest parametrem wiersza rejestru,
nie warunkiem zapisanym w kodzie.

## budowa/server/internal/dane/diagnostics_dziennik.go
Zawężenie idzie parametrem, nie sklejaniem tekstu. Jedno przygotowane
zapytanie obsługuje pięć zawężeń naraz, bo pusty parametr znaczy „nie
zawężaj". Wartości nigdy nie wchodzą do treści SQL, a pamięć podręczna
zapytań ma jedną pozycję zamiast trzydziestu dwóch.

## budowa/server/internal/session/petla.go

Pętla koordynator-wykonawca łączy w jednym bycie cztery obowiązki: przejście końca tury w kolejny
obieg koordynatora poprzez wybudzenie, licznik obiegów z nazwanym warunkiem zatrzymania przy braku
postępu, strumień wykonawcy widziany przez koordynatora oraz warunek ukończenia zadania, gdy tura
koordynatora zamyka się wynikiem przy braku okna wykonawczego w toku. Pętla nie powiela silnika
uruchamiania tury — samą turę koordynatora rozpoczyna warstwa rozmowy przez port UruchomienieObiegu.

Metoda ZakonczTure jest jedynym wejściem warstwy rozmowy do pętli i przyjmuje koniec tury każdego
okna, nie tylko wykonawczego. Koniec tury wykonawcy zdejmuje jego turę i wybudza koordynatora, skąd
wychodzi kolejny obieg. Koniec tury koordynatora zamyka bieg ukończeniem, gdy warunek metody
ukonczBieg jest spełniony; okno samodzielne nie porusza niczego, co nie jest usterką. Metoda zwraca
prawdę, gdy zgłoszenie poruszyło bieg: wybudziło koordynatora albo zamknęło go ukończeniem.

Warunek metody ukonczBieg jest podwójny i oba jego człony pętla mierzy, a nie zakłada: tura
koordynatora zamknęła się wynikiem kanału, przy czym warstwa rozmowy podaje wtedy powód wynikowy
wychodzący z tej samej trójki warunków co stan wiadomości ukończonej, oraz żadne okno wykonawcze
tego koordynatora nie prowadzi tury, czyli wykaz tur jego toru jest pusty. Bieg przed pierwszym
obiegiem się nie kończy, ponieważ nie ma czego kończyć, a zatrzymanie postawione oknu, które pętli
jeszcze nie prowadziło, wprowadziłoby fałszywy sygnał biegu nieistniejącego. Bieg już zatrzymany
zostaje przy swoim powodzie — ukończenie nie przykrywa przerwania.

## budowa/server/internal/models/wysylka.go
Kanał nierozpoznany kończy wyłącznie wywołanie Wyslij, którego dotyczy: użytkownik
dostaje fragment błędu, wywołujący błąd w wyniku, a sesja, okno i kolejne próby
pozostają czynne.

## budowa/server/internal/dane/dostep_korzenie.go
Korzeń jest listą, nie pojedynczym polem, tak samo jak katalog roboczy okna;
odpowiada zmiennej DANACO_MOST_KORZENIE mostu MCP, poza którą most nie wyjdzie.
Nadanie może korzenie punktu wyłącznie zawęzić. Korzeń nadania spoza obszaru
punktu byłby obietnicą dostępu, którego most i tak nie da — sprawdzenie leży
tu, żeby nie powstał wiersz wprowadzający w błąd.

Lista pusta w sprawdzeniu zawężenia znaczy „komplet korzeni punktu" i jest
poprawna zawsze. Punkt bez własnych korzeni nie ogranicza niczego — tak samo
jak most z pustą zmienną DANACO_MOST_KORZENIE.

## budowa/server/internal/dane/dostep_nadania.go
Nadanie wiąże okno komunikacji z punktem dostępu. Żyje per okno rozmowy, nie
per sesja i nie per platforma. Okno ma zbiór nadań — kolejność i oznaczenie
głównego niosą znaczenie. Okno bez nadań pracuje dalej, tylko niczego nie
widzi.

## budowa/server/internal/dane/agenci.go

TrybNakladki przyjmuje dwie wartości kontraktu: dołączenie do promptu systemowego jako
opcję domyślną albo zastąpienie go jako odstępstwo świadomie wybrane przez Operatora.

PoziomyPamieci: wycinek pusty jest jedynym zapisem wyłączenia pamięci w całości, nie ma
własnej piątej wartości poziomu, bo taka wartość dopuszczałaby stan sprzeczny z pozostałymi
czterema poziomami, którego nie dałoby się rozstrzygnąć.

LimitPodagentow: wartość wyjściowa piętnastu jest maksimum technicznym platformy.

ModulyZastosowania: wycinek pusty oznacza brak ograniczenia zastosowania, nie brak
odpowiedzi.

KodProjektu: pole puste zwraca bibliotekę w całości, ponieważ okno budowania eksperta musi
widzieć również ekspertów przypisanych do projektów, inaczej nie dałoby się ich poprawić.

Granica: liczba wszystkich wierszy spełniających warunki wraca osobno od granicy, żeby okno
wiedziało, ile pozycji zostało ucięte.

UstawPoziomyPamieci: rozróżnienie między pominięciem pola żądania a podaniem listy pustej
należy do warstwy wyższej, nie do repozytorium.

Wykaz zwracany przez zapytanie listaAgentow pomija archiwum, żeby ekspert odłożony nie
wisiał nadal na liście, co czyniłoby archiwizację znacznikiem bez skutku; archiwum ma własną
komendę odczytu.

## budowa/server/internal/dane/dostep_nadania_zapis.go
Reguły zbioru nadań: kolejność liczona jest od 1, brak wskazania dokłada
nadanie na koniec. Głównych nadań okna jest najwyżej jedno — pilnuje tego
indeks częściowy bazy, a repozytorium zdejmuje oznaczenie z poprzedniego,
zamiast zderzać się z więzem. Pierwsze nadanie okna zostaje główne z urzędu;
okno z nadaniami, ale bez głównego, nie miałoby punktu domyślnego.

## budowa/server/internal/dane/dostep_nadania_zbior.go
Reguły trzymane są osobno od poleceń zapisu, bo dotyczą całego zbioru nadań
okna, nie pojedynczego wiersza. Kolejność liczona jest od 1; brak wskazania
dokłada nadanie na koniec zbioru. Pierwsze nadanie okna zostaje główne
z urzędu — okno z nadaniami, ale bez głównego, nie miałoby punktu domyślnego.

## budowa/server/internal/models/zadanie_api.go
Pusty sekret pod istniejącym odwołaniem jest błędem tak samo jak referencja
sejfu bez wpiętego sejfu: kanał nie wysyła żądania, które i tak odbiłoby się
uwierzytelnieniem, a ciche zejście na zmienną środowiskową wzięłoby zmienną
o nazwie sejf:<byt>, której nikt nie ustawia. Sejf poświadczeń jest jedną
instancją nad jednym plikiem; montaż ustawia go raz na starcie, zanim serwer
zacznie obsługiwać żądania, więc odczyty idą już po zapisie. Klucz zapisany
komendą account.* jest tym samym, który kanał API odczytuje przy wysyłce,
ponieważ obie strony sięgają po ten sam sejf wpięty przy montażu.

## budowa/server/internal/session/petla_test.go

Sprawdziany warunku ukończenia biegu koordynator-wykonawca mierzą to, co pętla naprawdę dostaje przez
zgłoszenie ZakonczTure, i wykluczają dwie szkody: bieg, po którym nie da się maszynowo odróżnić pracy
skończonej od przerwanej, gdyż trzy powody zatrzymania mówią „przerwane" i żaden nie mówił „skończone
z wynikiem"; oraz ukończenie zamienione w bramkę akceptacji, czyli bieg, który po ukończeniu czeka na
potwierdzenie zamiast podjąć pracę zgłoszoną później. Rdzeń nie ma osobnego wyliczenia stanu tury —
tura jest wpisem w wykazie biegów warstwy rozmowy, a pętla dowiaduje się o jej końcu jednym zgłoszeniem.

## budowa/server/internal/dane/dostep_punkty.go
Punkt dostępu mówi, do czego model ma wgląd: maszyna udostępniona mostem MCP
albo katalog lokalny wskazanego urządzenia. Nie jest środowiskiem —
środowisko pozostaje profilem widoczności modułów w bocznej nawigacji. Nie
jest też katalogiem roboczym modelu: ten jest osobnym ustawieniem
(`katalog.roboczy.*`) i mówi, gdzie model zostawia własne pliki.

## budowa/server/internal/dane/dostep_punkty_wiersz.go
Adres mostu nie ma własnej kolumny — składa się go z użytkownika, hosta
i portu. Gdyby był osobnym polem, wiersz mógłby mieć adres sprzeczny z
własnymi danymi połączenia, a taki rozjazd nie ujawniłby się aż do nieudanego
uruchomienia.

## budowa/server/internal/dane/dostep_wyliczenia.go
Kontrakt nie deklaruje przy wyliczeniach AccessPointKind, AccessPointStatus
ani AccessMode pola bazy — inaczej niż przy AccountKind, które ma osobne
słowniki wartości bazy i wartości kontraktu. Własny przekład w warstwie
trwałości byłby drugim źródłem przekładu obok kontraktu współdzielonego.
Zbiory wyliczeniowe tego pliku są więc sprawdzeniem przynależności, nie
przekładem: pilnują, żeby do kolumny nie trafiła wartość spoza kontraktu.

## budowa/server/internal/dane/agenci_zapis.go

Uprawnienie jest konfiguracją możliwości, nie bramą: nowo założony ekspert pracuje od razu,
a interfejs uprawnień pokazuje mu stan wyjściowy zamiast pustego formularza.

widocznoscKolumny: wartość spoza katalogu widoczności odbija się o warunek CHECK kolumny
agent.widocznosc bazy danych, więc drugi katalog dopuszczalnych wartości nie powstaje w kodzie.

## budowa/server/internal/dane/agent_archiwum.go

Wybór między usunięciem a archiwizacją eksperta należy do Operatora, nie do rdzenia, więc
rdzeń nie przekierowuje jednej komendy na drugą; usunięcie zostaje jedyną drogą utraty
definicji eksperta, wzorem kosza sesji.

Stan czynności eksperta wraca po przywróceniu taki, jaki był przed archiwizacją, zamiast
zakładać, że każdy zarchiwizowany był czynny — ekspert wyłączony przed archiwizacją
wróciłby inaczej, niż go odkładano.

Migawek historii wersji archiwizacja nie mnoży: zapis nie rusza licznika wersji eksperta,
więc wyzwalacz bazy uzupełnia migawkę wersji bieżącej zamiast zakładać nową. Odłożenie
eksperta na półkę nie jest zmianą jego tożsamości.

## budowa/server/internal/models/zapytanie.go
Kanał nie rozstrzyga konfiguracji wywołania sam — dostaje zasięgi po to, by móc
je przekazać dalej i odnotować w prowenancji.

Historia jest strukturą wewnętrzną pakietu, nie kontraktem: adapter kanału
składa z niej pamięć wywołania, a wypełnia ją warstwa rozmowy z dziennika tury.
Kanał bezstanowy składa z niej pamięć wywołania przy każdym wywołaniu; kanał
z własną pamięcią może ją pominąć.

Rozdzielenie ról obrazu wejściowego na materiał i maskę jest konieczne, bo dwie
czynności warsztatu fotografii — uzupełnienie ubytku i rozszerzenie kadru —
wysyłają obie role naraz, a punkty końcowe dostawcy przyjmują je osobnymi
polami.

Pole ObrazWejsciowy istnieje, bo bez niego kanał obrazowy umiał wyłącznie
wygenerować obraz z samego polecenia tekstowego: cztery czynności warsztatu
fotografii modułu Design — powiększenie, odcięcie tła, uzupełnienie ubytku
i rozszerzenie kadru — nie miały jak dojść do wariantu neuronowego, choć mają
go lepszy od rachunku na pikselach. Bajty jadą wyłącznie wprost, bez wariantu
adresowego, jaki ma treść obrazu w drugą stronę: odpowiedź dostawcy bywa
odsyłaczem, bo to dostawca trzyma plik, a materiał wejściowy leży w magazynie
rdzenia, do którego punkt końcowy dostawcy nie ma jak sięgnąć — adres przyjęty
tutaj byłby adresem, którego druga strona nie odczyta.

Nazwa pliku obrazu wejściowego pozwala punktowi końcowemu rozpoznać format,
gdy nie dostanie typu treści.

ObrazyWejsciowe: kanał obrazowy bierze pierwszy materiał i pierwszą maskę
z listy, bo punkty końcowe edycji przyjmują po jednym elemencie z każdej roli.

Pułap kosztu tury pochodzi z nastawy pulap_kosztu_usd; zerowa wartość oznacza,
że kanał nie dopisuje przełącznika ograniczającego budżet wywołania.

Katalog roboczy sesji ustala się z konfiguracji operatora, kluczami obszaru
katalogu roboczego.

Konfiguracja mostów MCP sesji jedzie osobnymi przełącznikami obok konfiguracji
mostów okna, więc nadania okna i wiązania sesji nie przykrywają się nawzajem.

Plik ustawień sesji składa się z obszarów narzędzi i uprawnień konfiguracji
sesji.

Środowisko procesu kanału dokładają zmienne wyliczone z obszarów środowiska
i dostawcy w konfiguracji sesji.

Konto wybrane dla wywołania: gdy zapytanie nie wskazuje konta wprost, obowiązuje
konto powiązane z wierszem kanału w rejestrze. Powiązanie kanału z kontem staje
się dzięki temu widoczne w wywołaniu nawet bez konfiguracji sesji; brak obu
źródeł oznacza, że tożsamość bierze się z otoczenia procesu.

## budowa/server/internal/session/proces.go

Proces nie powstaje w tym pakiecie ani nie jest stąd uruchamiany. Startuje go warstwa kanału własną
drogą, a sesja obejmuje proces już biegnący przez metodę Przejmij rejestru procesów. Do sesji należy
wyłącznie to, czego kanał nie umie: objęcie całego drzewa potomstwa jednym uchwytem systemowym, czyli
Job Object na Windows albo grupa procesów na systemach uniksowych. Ubicie okna kończy zatem także
wnuki, bez narzędzia taskkill i bez innej zależności od narzędzi systemu.

Bez doglądu wykonywanego metodą dogladaj pole zakonczony miałoby jednego pisarza, metodę Ubij, a do
Ubij prowadzą wyłącznie dwie drogi: zatrzymanie okna i przejęcie procesu następnej tury. Proces, który
kończy się sam, nie wyzwala żadnej z nich, więc okno meldowałoby stan biegnący od samoistnego wyjścia
aż do najbliższej tury albo zamknięcia okna, a przez ten czas wisiałby uchwyt zadania i uchwyt procesu
systemowego.

Metoda zwalniająca uchwyty ma teraz dwóch wołających, Ubij i dogląd, a zwolnienie drzewa procesów nie
jest współbieżnie idempotentne: dwa zamknięcia tego samego uchwytu zamykają uchwyt, który system
zdążył już nadać ponownie czemu innemu. Bezpieczeństwo stoi wyłącznie na odczycie i zapisie pola
zakonczony pod tą samą blokadą, co w metodzie Ubij: kto zastanie wartość fałszywą, ten jeden przechodzi
dalej. Tego nie wolno uprościć do sprawdzenia stanu przed działaniem, ponieważ byłoby to sprawdzenie
i działanie rozdzielone w czasie, przy którym obaj wołający mogliby wejść równocześnie.
## budowa/server/internal/dane/poczta.go
Repozytorium nie mówi żadnym protokołem: odbiór i wysyłka listów należą do warstwy poczty, tu leży wyłącznie trwałość. Hasła tu nie ma: kolumna z odwołaniem do poświadczeń niesie wskazanie na sejf poświadczeń, a repozytorium nigdy nie widzi sekretu w jawnej postaci.

## budowa/server/internal/dane/dostep_punkty_zapis.go
Punkt bez swoich korzeni byłby punktem, który obiecuje dostęp do całego
systemu plików maszyny. Kolumna poświadczenia niesie wyłącznie nazwę wpisu
w magazynie sekretów — repozytorium nie ma metody zapisującej treść klucza.

Warstwa wyższa odróżnia brak urządzenia bieżącego od innych błędów i zamienia
go na odmowę z powodem, nie na awarię.

Schemat żąda urządzenia od katalogu lokalnego, a klucz obcy żąda, żeby
wskazane urządzenie istniało. Oba więzy zgłaszają się dopiero jako awaria
zapisu, więc wybór urządzenia zapada przed poleceniem — inaczej zamiast
powodu odmowy wraca błąd wewnętrzny. Wskazanie podane sprawdzamy istnieniem
wiersza. Katalog lokalny bez wskazania osadzamy na maszynie, na której działa
rdzeń: katalog wybrany oknem powłoki leży z definicji tam, a kolumna
oznaczająca maszynę bieżącą tę maszynę wskazuje. Punkt mostowy bez wskazania
zostaje bez urządzenia — most jest wpisem platformy i urządzenia nie wymaga.

## budowa/server/internal/dane/extension_integracje.go
Poświadczenie jest w tym pliku wyłącznie odwołaniem — kluczem jawnym warstwy
sekretów. Hasła, tokenu ani klucza API ta warstwa nie widzi i nie ma jak
zobaczyć.

## budowa/server/internal/dane/agent_warstwy.go

Warstwa promptu jest stanem, nie zdarzeniem: klucz główny tabeli warstw stoi na parze
identyfikatora eksperta i nazwy warstwy, więc powtórzony zapis tej samej warstwy nadpisuje
treść zamiast dokładać drugi wiersz; UstawWarstwe wywołane dwa razy z tą samą treścią
zostawia bazę w tym samym stanie.

Nazwy warstw i tryby podania stoją wyłącznie w warunkach CHECK bazy, wartościami kontraktu
wprost; drugiego katalogu dopuszczonych wartości repozytorium nie prowadzi.

Warstwa o nazwie spoza katalogu ląduje na końcu porządku krytyczności.

Tryb nałożenia instrukcji dotyczy eksperta jako całości, nie pojedynczej warstwy: prompt
globalny albo obowiązuje w całości, albo nie — zastąpienie częściowe, samą warstwą profilu,
nie ma znaczenia, więc kolumna trybu stoi przy ekspercie, nie przy każdej z warstw osobno.

## budowa/server/internal/session/przejecie.go

Kanał modelu startuje własny proces, ponieważ tylko on zna wiersz poleceń, rotację kont i kształt
strumienia. Ubicie drzewa procesów pozostaje jednak w jednym miejscu — w tym pliku; bez przejęcia po
ubiciu okna zostałyby wnuki procesu kanału, takie jak serwery narzędziowe czy powłoki uruchomione przez
model. Warstwa kanału używa obu części naraz: nadaje atrybuty uruchomienia przed startem procesu,
przejmuje drzewo zaraz po uruchomieniu, zwalnia je odroczonym wywołaniem i ubija przy zamknięciu okna.
Pominięcie atrybutów uruchomienia nie wywraca kanału — na Windows przejęcie zadziała mimo to, a na
systemach uniksowych ubicie obejmie sam proces zamiast całej grupy, ponieważ brak elementu opcjonalnego
nie blokuje uruchomienia.

Czekanie nie rusza uchwytów oddawanych przez zwolnienie, więc wolno je prowadzić równolegle z ubiciem:
obserwator obudzi się wtedy, gdy ubicie zrobi swoje. Na Windows czekanie oznacza oczekiwanie na uchwyt
procesu, a na systemach uniksowych odpytywanie sygnałem zerowym.

## budowa/server/internal/nadajnik/nadajnik.go
Pakiet wysyła dokładnie dwa listy: potwierdzenie adresu przy rejestracji
i drogę odzyskania konta, żadnych innych. Różni się od modułu poczty, który
jest klientem skrzynki Operatora: czyta jego listy, wysyła w jego imieniu
i odkłada kopię w jego folderze wysłanych. Nadajnik nadaje w imieniu platformy,
do Operatora, i kopii nigdzie nie odkłada — list systemowy w folderze wysłane
Operatora byłby śladem czynności, której on nie wykonał. Stąd osobne
poświadczenie i osobny host: konto nadawcze platformy nie jest skrzynką
Operatora i nie wolno ich mieszać. Gdyby aplikacja pisała jego kontem, utrata
dostępu do skrzynki odcinałaby drogę odzyskania konta, dokładnie wtedy, gdy
jest potrzebna.

List niesie drogę potwierdzenia, nie hasło: platforma nie zna hasła w postaci
jawnej i nigdy go nie odsyła. Wysyłka SMTP jedzie biblioteką standardową, bez
dodatkowej zależności.

Sekret konta nadawczego przychodzi z sejfu poświadczeń, a pakiet traktuje go
jako nieprzezroczysty napis: nie zapisuje go, nie loguje i nie umieszcza
w treści listu.

Brak konta nadawczego nie jest usterką do zgłoszenia w połowie rejestracji —
to stan, o którym warstwa wyżej musi wiedzieć, zanim założy konto i obieca
list.

Uwierzytelnienie w rozmowie SMTP jest warunkowe: serwer dostawcy zawsze go
żąda, ale przekaźnik na tej samej maszynie często nie ogłasza go wcale —
wpychanie mu wtedy poświadczenia kończy się odmową przy komendzie, która bez
uwierzytelnienia by przeszła.
## budowa/server/internal/dane/poczta_skrzynki.go
Ta sama tabela skrzynka_pocztowa niesie dwa różne pytania: repozytorium poczty czyta ją pod kątem obserwatora poczty (skrzynki czynne, takt odpytywania), a to repozytorium pod kątem rodziny mail (podpięcie, wykaz, domyślna). Drugiej tabeli na skrzynkę nie ma: gdyby była, automatyka wyzwalana listem nie widziałaby skrzynki podpiętej komendą. Hasła tu nie ma: kolumna z odwołaniem do sejfu niesie wskazanie na wpis sejfu poświadczeń, a kolumny na sam sekret schemat nie zna, więc żaden odczyt tego repozytorium nie ma jak go wynieść. Zdjęcie domyślności z pozostałych skrzynek i nadanie jej nowej idą jedną transakcją, bo indeks częściowy schematu dopuszcza jedną domyślną skrzynkę naraz.

## budowa/server/internal/dane/jakosc.go
Zapis jest wymianą, nie dokładaniem. Kontrakt odpowiedzi kontroli jakości nie
niesie pola przyrostowego, ani identyfikatora kontroli poprzedniej — oddaje
płaski wykaz zastrzeżeń bieżącego stanu panelu. Każde wywołanie kontroli
liczy niezgodności na nowo z treści panelu i zastępuje poprzedni wykaz tej
samej kontroli, tym samym wzorcem co zapis kroków automatyzacji. Dopisywanie
dawałoby narastającą listę powtórzeń tej samej usterki przy każdej kolejnej
kontroli tego samego panelu.

## budowa/server/internal/dane/jakosc_eksport.go
Ten plik zawiera wyłącznie zapis i odczyt śladu eksportu; typ, interfejs
i konstruktor repozytorium tłumaczeń leżą w innym pliku. Rdzeń nie ma
magazynu plików binarnych, więc zlecenie eksportu nie wytwarza pliku
i odnośnik do pliku zostaje pusty, dopóki plik realnie nie powstał poza
rdzeniem. Zapisywany jest wyłącznie ślad zlecenia eksportu: w jakim formacie
i kiedy. Czas jest liczbą milisekund epoki, tym samym wzorem co w pozostałych
miejscach warstwy trwałości.

Panel, którego nie ma, wraca jako błąd braku wiersza — cicha zgoda na eksport
bytu, którego nie ma, byłaby potwierdzeniem czynności, która się nie odbyła.

## budowa/server/internal/dane/jakosc_mowa.go
Zasila komendę syntezy mowy; plik jest osobny od pliku niezgodności jakości
i od pliku śladów eksportu. Rdzeń nie syntezuje mowy — ten sam brak jak
w module Assistant. Odnośnik do nagrania niesie odwołanie do pliku
dostarczonego z zewnątrz, jeśli kiedykolwiek powstanie; metoda zapisu nie
dorabia mu wartości domyślnej — brak zostaje pusty, bo rdzeń nie syntezuje,
a zmyślona ścieżka byłaby obietnicą bez pokrycia.

## budowa/server/internal/narzedzia/adres.go
Port czyta pakiet konfiguracja, ten sam, który ustala port nasłuchu procesu
rdzenia, więc zmiana portu przez Operatora przestawia obie strony naraz.
Odczyt zmiennej środowiska po nazwie dosłownej należy wyłącznie do tamtego
pakietu, a serwer narzędzi dziedziczy środowisko po procesie modelu, który
dziedziczy je po rdzeniu.

Pętla zwrotna jest wyborem, nie skrótem: serwer narzędzi stoi zawsze na tej
samej maszynie co proces modelu, który go uruchomił, a proces modelu stoi
przy rdzeniu, który go zrodził. Adres inny niż pętla zwrotna byłby wtedy
zgadywaniem — Operator wskazuje go przełącznikiem, gdy układ jest inny.

## budowa/server/internal/session/rejestr_procesow.go

Rejestr nie uruchamia procesów i nie zna drogi ich uruchomienia. Proces tury startuje warstwa kanału,
a do rejestru trafia przez metodę Przejmij, objęty uchwytem drzewa i gotowy do zatrzymania. Jest to
jedyna droga wpisu do tego rejestru: bez przejęcia zamknięcie okna nie zatrzymałoby tego, co model
uruchomił. Przejęcie zakłada Job Object na Windows albo grupę procesów na systemach uniksowych, dzięki
czemu ubicie okna kończy także wnuki procesu, bez narzędzia taskkill i bez zależności od narzędzi
systemu. Niepowodzenie przejęcia nie przerywa tury: proces biegnie i odpowiada, tylko jego potomstwo
nie jest objęte uchwytem.

Dogląd uruchomiony po wstawieniu wpisu do mapy jest jedyną drogą obserwacji, więc pokrywa każdy wpis,
i żaden proces kończący się sam nie zostaje w rejestrze jako biegnący. Gorutyna doglądu nie sięga po
blokadę rejestru, więc jej start pod zamkiem niczego nie blokuje.

## budowa/server/internal/narzedzia/ekspert_definicja.go
Serwer narzędzi jest dla rdzenia zwykłym urządzeniem: ta sama koperta, to samo
gniazdo, co okno interfejsu. Drugiego wejścia do danych eksperta tu nie ma
i nie powstaje — sięgnięcie po sterownik bazy z procesu modelu byłoby
obejściem rdzenia, a nie skrótem.

Kontrakt nie ma osobnej komendy do odczytu jednego eksperta, a żądanie wykazu
agentów nie filtruje po identyfikatorze. Wykaz jest więc pobierany w całości
i dopasowywany po identyfikatorze eksperta; gdyby komenda odczytu pojedynczego
eksperta kiedyś powstała, funkcja OdczytajEksperta jest jedynym miejscem do
zmiany.

Gniazdo do rdzenia zestawia się przy pierwszym użyciu: proces modelu uruchamia
serwery MCP na starcie rozmowy, a rdzeń bywa wtedy jeszcze niegotowy.
Definicji nie da się więc mieć w chwili startu procesu, dlatego czyta się ją
przy pierwszym żądaniu tools/list.

Pola Umiejetnosci i Konektory typu DefinicjaEksperta odpowiadają polom kontraktu
niosącym umiejętności i konektory eksperta, przeniesionym bez zmiany znaczenia;
reszta eksperta — warstwy, model, uprawnienia — do doboru narzędzi nie należy
i nie jest tu kopiowana.

Kody wraca jedną listą, bo dobór narzędzi nie rozróżnia pochodzenia kodu:
rozpoznanie idzie po tym, czy kod nazywa narzędzie albo grupę, a nie po tym,
w którym polu eksperta go zapisano.

Rozstrzygnięcie, co zrobić z niewiedzą o eksperckim wyposażeniu, należy do
składania wykazu, nie do OdczytajEksperta — tutaj jest wyłącznie odczyt
i wyłącznie prawda o tym, co rdzeń powiedział.

Kod eksperta nieznany rdzeniowi jest faktem, nie pustką: ekspert bywa kasowany
niezależnie od okien, w których pracował, więc okno może nieść kod, którego
biblioteka już nie ma. Taki stan wraca błędem mówiącym, ilu ekspertów rdzeń
zna, żeby stan nieznaleziony był odróżnialny od stanu, w którym rdzeń oddał
wykaz pusty.
## budowa/server/internal/dane/przegladarka_notatki.go
Kolumna wskazująca źródło zewnętrzne jest wartością tekstową, nie więzem obcym: notatka może dotyczyć całej strony, nie tylko jednego zebranego źródła, więc kolumna jest dopuszczalnie pusta bez odwołania referencyjnego. Treść notatki jest krótkim tekstem wprost, nie odwołaniem do pliku, w odróżnieniu od migawki strony, która trzyma treść obszerną osobno — notatka Operatora nią nie jest, więc kolumna niesie treść wprost.

## budowa/server/internal/dane/historia_retencja.go
Osobny plik obok pliku pozycji historii dzieli dwie odpowiedzialności: tam
żyje pozycja historii — wiersz wiadomości czytany i kasowany na wskazanie
operatora; tu żyje nastawa, która kasuje sama, bez wskazania. Zasady się nie
sumują. Obowiązuje jedna — najbliższa oknu (okno, potem sesja, potem
globalna). Suma dawałaby wynik, którego operator nie przewidziałby z żadnego
pojedynczego ekranu.

Byt zakresu, którego nie ma, nie da się później przyciąć niczym, więc zasada
zapisana dla nieistniejącego okna albo sesji byłaby nastawą bez skutku.

Rozstrzygnięcie zasady idzie po kolejności okno, sesja, globalna, z granicą
jednego wiersza, więc wiersz pusty zakresu węższego przesłaniałby zasadę
szerszą, a jedyna komenda ustawiania retencji nie ma czym takiego wiersza
skasować. Zdjęcie wiersza sprawia, że brak wiersza znaczy „ten zakres nic nie
postanawia", a nie „trzymaj zero" — pytanie o zasadę spada wtedy na zakres
szerszy, dokładnie jak przed pierwszym zapisem.

Zakres globalny bytu nie wskazuje i jest zawsze prawdziwy — tnie wszystkie
okna. Fałsz nie jest błędem odczytu: to stan, który woła o odmowę po stronie
komendy, bo zasada zapisana dla nieistniejącego okna albo sesji nigdy
niczego nie przytnie.

## budowa/server/internal/dane/agent_wersje.go

Migawki historii zakłada baza, nie repozytorium: wiersz historii powstaje wyzwalaczem, bo
tożsamość eksperta zapisują trzy różne drogi kodu — założenie, zmiana i przywrócenie.
Repozytorium historii wyłącznie czyta i przywraca; nie ma osobnej czynności zapisu wersji,
bo taka czynność pozwalałaby historię ominąć.

Różnicę względem poprzedniej wersji wyliczamy przy odczycie zamiast trzymać ją w bazie: pola
zmienione wychodzą z porównania migawki z jej poprzedniczką, więc odpowiedź nie rozjedzie się
z treścią, którą przywraca operacja przywrócenia.

## budowa/server/internal/narzedzia/ekspert_wykaz.go
Umiejętności i konektory eksperta są w kontrakcie listami napisów bez
narzuconego słownika, wpisywanych ręcznie. Kod rozpoznaje się dwiema drogami
sprawdzalnymi wprost wobec kontraktu: nazwa narzędzia wskazuje jedną pozycję
wykazu, nazwa grupy wskazuje wszystkie pozycje tego obszaru. Grupa jest
jednostką doboru właśnie po to, żeby dobór dawał się wykonać jednym
wskazaniem zamiast wieloma osobnymi.

Kod, którego nie da się rozpoznać żadną z dwóch dróg, nie jest połykany po
cichu: wraca osobną listą kodów nierozpoznanych i idzie stamtąd do dziennika
oraz do licznika.

Zawężenie wchodzi wtedy i tylko wtedy, gdy jest czym zawęzić, czyli gdy
rozpoznano co najmniej jeden kod. Każdy inny przypadek — rdzeń milczy, kodu
nie ma, żadnego kodu nie rozpoznano — zostawia wykaz okna w całości i melduje
powód.

Wykaz okna w całości wybrano zamiast wykazu pustego, bo usterka jednej drogi
nie zabiera modelowi wszystkich narzędzi na całą turę: okno z ekspertem,
którego rdzeń chwilowo nie potwierdził, pracuje dalej. Pełny wykaz kosztuje
jednak istotnie więcej żetonów niż wykaz zawężony, dlatego koszt jest mierzony
i meldowany przy każdym takim przypadku. Wykaz pusty odpada, bo model bez ani
jednego narzędzia nie mówi „nie znam eksperta", tylko po prostu nie działa.

Nazwa przełącznika dołożeń i jego rozdzielnik pochodzą ze wspólnego pakietu
wstrzykiwania, nie z literałów zapisanych lokalnie: wiersz uruchomienia składa
strona rdzenia, a czyta go serwer narzędzi z przetwarzaniem kończącym proces
na błędzie, więc rozjazd nazwy o jeden znak ubiłby cały proces zamiast
zawężyć jedynie turę o kilka pozycji.

Dołożenie w ZDolozeniami dokłada i nigdy nie zawęża: dorzuca narzędzie do
pracy poza definicją eksperta i bez jej ruszania. Pozycje dołożone idą na
koniec listy, żeby kolejność wykazu została kolejnością kontraktu, a dołożenie
było widoczne jako dołożenie. Nazwa nierozpoznana dopisuje się do listy
nierozpoznanych także wtedy, gdy wykaz nie był zawężony i dokładać nie było
czego, bo nazwa nienazywająca ani narzędzia, ani grupy nie zadziała nigdy.

WykazBezEksperta jest osobną funkcją, a nie wynikiem składanym w miejscu
wywołania, bo to jest przypadek, o którym najłatwiej zapomnieć powiedzieć —
powód wchodzi tu zawsze, więc nie da się złożyć niewiedzy bez jej opisania.

## budowa/server/internal/narzedzia/grupa.go
Kontrakt stanowi notację obszar-nazwa i sam się nią posługuje przy rozpoznaniu
zdarzeń nieznanych, więc grupa jest tu wyprowadzana z nazwy komendy, a nie
zapisana osobnym wykazem; wykaz własny narzędzie-grupa byłby drugą prawdą,
która rozjedzie się przy pierwszym narzędziu dopisanym do kontraktu.

Nazwa grupy jest kodem obszaru, nie zdaniem opisującym, do czego dany obszar
służy: takie zdania mieszkają w sekcji obszarów kontraktu, ale generator nie
emituje ich do wygenerowanego kodu Go. Opis zmyślony byłby zdaniem o
kontrakcie, którego kontrakt nie mówi, więc do czasu udostępnienia opisów
grupa niesie sam kod obszaru.

Rozstrzygnięcie grupy dla komendy bez separatora — grupa równa całej nazwie —
jest zgodne z rozpoznaniem zdarzeń nieznanych w pakiecie współdzielonym, żeby
dwa różne rozstrzygnięcia tej samej notacji nie rozjechały się.

Kolejność alfabetyczna nazw grup, a nie kolejność kontraktu, wynika z tego, że
grupy są gałęziami drzewa wyboru w interfejsie i gałąź szuka się okiem po
nazwie; kolejność kontraktu zostaje tam, gdzie ma znaczenie — wewnątrz wykazu
pozycji.
## budowa/server/internal/dane/przegladarka_zestawy.go
Skład zestawu i wątku nie jest osobną kolumną: źródło wskazuje swój zestaw kolumną grupującą, notatka swój wątek kolumną wątku, więc skład jest zapytaniem po tej kolumnie. Druga lista, czyli osobny wykaz identyfikatorów zapisany przy zestawie, rozjechałaby się z pierwszą przy pierwszym usunięciu źródła, a rozjazd nie byłby widoczny, bo obie listy wyglądałyby wiarygodnie.

## budowa/server/internal/dane/agent_wtyczki.go

Wtyczka nie jest konektorem. Konektor jest drogą do usługi: wskazuje most z katalogu punktów
dostępu, z którego rdzeń składa listę serwerów MCP podawaną przełącznikiem konfiguracji.
Wtyczka jest katalogiem rozszerzeń powłoki — nie ma adresu ani poświadczenia, ma nazwę,
źródło i wersję, a program dostaje ją osobnym przełącznikiem katalogu wtyczek. Stąd osobna
tabela i osobne wejście repozytorium dla obu pojęć.

Kod wtyczki nadaje repozytorium, w odróżnieniu od konektora, któremu kod nadaje wołający,
ponieważ struktura konektora wchodzi w całości od wołającego, a wejście wtyczki bierze
wyłącznie nazwę, źródło i wersję. Kod składa się dokładnie tak samo jak dla konektora:
przedrostek, licznik w podstawie trzydziestej szóstej i ośmiobajtowa część losowa ze źródła
kryptograficznego.
## budowa/server/internal/dane/przegladarka_zrodla.go
Tabela zrodlo_przegladania nie jest tabelą zrodlo_badania modułu badawczego: źródło przeglądania jest odciskiem strony zebranym w toku przeglądania i zawsze powiązanym z oknem operacyjnym, a źródło badawcze ocenia wiarygodność zasobu i niesie inny kształt danych. Różne kształty i różne cykle życia uzasadniają osobną tabelę zamiast współdzielenia jednej struktury.

## budowa/server/internal/narzedzia/licznik.go
Sama liczba pozycji nie mówi Operatorowi nic; liczba bajtów i rząd żetonów okna
kontekstu mówi wszystko, bo limit, w który się uderza, jest limitem okna
kontekstu liczonym w żetonach, nie w narzędziach. Licznik podający wyłącznie
pozycje kazałby Operatorowi przeliczać je w głowie na koszt, czyli o
przekroczeniu dowiedziałby się dopiero z cichej degradacji.

Bajty mają być bajtami tej samej odpowiedzi, którą dostanie model, inaczej
pomiar jest oszacowaniem podanym jako pomiar. Trzy pola, z których protokół
składa odpowiedź tools/list, stoją w jednym miejscu w tym pliku; warstwa
protokołu bierze je stąd, zamiast składać drugi raz po swojemu. Pakiet nadal
nie zna ramki JSON-RPC ani metod, zna wyłącznie kształt danych jednej pozycji,
i to jest cena za to, że pomiar nie kłamie.

Licznik oddaje bajty, bo bajty umie policzyć dokładnie. Przelicznika na żetony
tu nie ma: zależy od tokenizatora kanału modelu, którego ten proces nie zna,
a liczba podana jako dokładna, wyprowadzona z założenia, byłaby drugą prawdą.
Rząd wielkości podaje się w zdaniu dziennika, nie w polu struktury.

## budowa/server/internal/dane/kanaly.go
Baza przechowuje wyłącznie odwołanie do danych dostępowych — nazwę wpisu
w magazynie sekretów, nigdy klucza. Repozytorium nie ma żadnej metody
zapisującej treść sekretu, a parametry kanału są sprawdzane jako poprawny
dokument JSON.

## budowa/server/internal/dane/karty_sesji.go
Repozytorium ma metodę zapewniającą istnienie karty: warstwa wyższa zapisuje
sesję, a karta ma powstać po drodze, nie zablokować zapisu.

## budowa/server/internal/dane/katalog_definicji.go
Pozycja katalogu odpowiada strukturze definicji ustawienia kontraktu —
klient buduje z niej pole formularza i nie zna ani jednego klucza z osobna.
Zbiory poboczne (opcje, zasięgi, osie) czytane są trzema zapytaniami zbiorczo
i dokładane do pozycji w pamięci; zapytania per pozycja nie ma.

## budowa/server/internal/dane/katalog_ustawien.go
Katalog jest sterowany danymi: nowa pozycja okna konfiguracji to nowy
wiersz, nie nowa gałąź w kodzie. Wzorcem jest rejestr kanałów modelu
i katalog akcji — repozytorium wyłącznie czyta wiersze, a rozstrzyganie
wartości należy do warstwy konfiguracji. Brak wiersza w katalogu nie jest
awarią: rezolwer schodzi wtedy na rejestr wbudowany rdzenia i pracuje dalej.

Rodzaje liczbowe, logiczne i złożone idą surowo, jeżeli są poprawnym
zapisem JSON; wszystko pozostałe idzie napisem. Funkcja nigdy nie zawodzi —
wartość nieczytelna trafia do kontraktu jako napis, nie jako błąd.

## budowa/server/internal/dane/kolejki.go
Zmiana stanu dotyka stanu i dziennika akcji, więc idzie w transakcji.
Zlecenia kolejki obsługuje osobny plik repozytorium pozycji kolejki.

Rachunek liczby kolejek czynnych stoi w tym repozytorium, bo tabelę kolejek
prowadzi ono i drugiego czytelnika mieć nie będzie; ciało metody leży
w pliku warstwy mobilnej.

## budowa/server/internal/dane/akcje.go

Filtr do zasięgu celowo nie stoi w zapytaniu SQL: wybór akcji jednego bytu poziomu rozstrzyga
wyłącznie rejestr akcji w rdzeniu na wierszach już odczytanych, żeby bliźniaczy filtr w SQL
nie stał się drugą, osobną implementacją tej samej reguły, gotową rozejść się z pierwszą po
cichu.
## budowa/server/internal/dane/przekazanie_okna.go
Interfejs kontraktu obszaru window.* deklaruje w całości wyłącznie ten plik, wraz z metodami, które implementują pozostałe pliki obszaru: więź koordynator-wykonawca oraz dziennik akcji. Interfejs rozdzielony na kilka plików byłby kilkoma prawdami o jednym kontrakcie. Identyfikator pozycji kolejki nie jest zakładany własnym zapisem w tym repozytorium: pozycję kolejki zakłada jedyny silnik kolejek, a adapter rdzenia wypełnia to pole gotowym identyfikatorem po założeniu pozycji, w tej samej transakcji co zapis zlecenia. Komplet kontekstu jest przechowywany w całości jako surowy zapis JSON: warstwa danych go nie interpretuje ani nie rozbiera na pola, tylko przechowuje i oddaje.

## budowa/server/internal/narzedzia/polaczenie.go
Serwer narzędzi jest dla rdzenia zwykłym urządzeniem, bez drugiego wejścia do
rdzenia. Połączenie jest leniwe: proces modelu uruchamia serwery MCP na
starcie rozmowy, więc odmowa startu przy niedostępnym rdzeniu zabrałaby
modelowi wszystkie narzędzia na całą turę. Serwer wstaje zawsze, a gniazdo
zestawia się przy pierwszym wywołaniu; nieudane zestawienie wraca do modelu
treścią błędu i nie przeszkadza próbie następnej.

Serwer narzędzi przedstawia się przy nawiązaniu połączenia, nie powitaniem,
bo powitania nie wysyła: jest klientem wołającym komendy, a nie oknem
interfejsu. Adres niesie rodzaj klienta, rolę okna i samo okno, dzięki czemu
rdzeń wie, czy po drugiej stronie stoi klawiatura Operatora czy zwykły model
roboczy; bez tego rdzeń widziałby wyłącznie identyfikator gniazda. To
przedstawienie nie jest uprawnieniem — rdzeń niczego na nim nie warunkuje,
wpisuje je do pola opisowego zdarzenia. Zasięg narzędzi rozstrzyga nadal
wyłącznie przełącznik uruchomieniowy czytany z wpisu MCP ułożonego przez
rdzeń, a nie ten napis tożsamości.

Adres nieczytelny w zAdresemTozsamosci zostaje adresem dotychczasowym:
gniazdo bez tożsamości jest gorsze od gniazda z tożsamością, ale
nieporównanie lepsze od braku narzędzi przez cały czas życia procesu modelu.

Gniazdo niesie także rozgłoszenia rdzenia — zdarzenia i fragmenty strumienia
innych okien — dlatego Wykonaj rozpoznaje odpowiedź po trzech rzeczach naraz:
identyfikatorze żądania, nazwie komendy i obecności pola stanu. Rozgłoszenie
nie ma stanu i nosi nazwę zdarzenia, więc nie da się go wziąć za odpowiedź.

## budowa/server/internal/dane/konta_wiersz.go
Kolumna poświadczenia dochodzi tutaj wyłącznie jako znacznik obecności
policzony w zapytaniu — struktura nie ma pola na jej treść, więc odczyt
katalogu nie ma czym wynieść odwołania.

## budowa/server/internal/dane/konta_domyslne.go
Repozytorium ma jedynie zdjąć oznaczenie domyślności z poprzedniego konta
i nadać je nowemu w jednej transakcji, żeby indeks nigdy nie zobaczył dwóch
kont domyślnych naraz.

## budowa/server/internal/dane/konta_powiazanie.go
Wiązaniem jest kolumna kanału modelu wskazująca konto, z kasowaniem
kaskadowym do wartości pustej, i nic poza nią — konto i kanał to dwa różne
byty. Kanał jest definicją rozmowy z modelem: jak wołać, jakim modelem,
z jakimi parametrami. Konto jest profilem uwierzytelnienia. Jeden kanał
wskazuje konto preferowane, jedno konto może obsługiwać wiele kanałów,
a pula rotacji bierze konta tego samego rodzaju — dlatego kanał pracuje
dalej także wtedy, gdy jego konto preferowane wyczerpało limit.

Z tego wynika sposób usuwania: skasowanie konta odłącza kanały, ale ich nie
kasuje. Kontrakt oddaje to wykazem odłączonych kanałów w odpowiedzi
usunięcia konta, więc repozytorium musi odczytać wykaz kanałów przed
skasowaniem wiersza — po skasowaniu wiązania już nie ma.
## budowa/server/internal/dane/przekazanie_okna_akcje.go
Katalog akcji mówi, jakie akcje istnieją; ten dziennik mówi, kiedy i z jakim skutkiem konkretne okno je wykonało, na wzór dziennika akcji kolejki: przejrzystość zamiast bramy. Typ repozytorium i konstruktor deklaruje plik sąsiedni tego samego obszaru; ten plik dokłada wyłącznie metody dziennika akcji. Parametry i wynik są surowym zapisem, nierozbieranym, bo kształt obu pól zależy od konkretnej akcji z katalogu, którego warstwa danych akcji nie zna, podobnie jak komplet kontekstu przekazania niesie treść bez rozbioru w zapytaniu. Pole niesie identyfikator zewnętrzny okna, ten sam, którym okno wychodzi kontraktem na warstwę wyższą, bo warstwa wyższa nie zna wewnętrznych kluczy liczbowych.

## budowa/server/internal/session/ubicie_unix.go

Proces okna zakłada własną grupę procesów, a potomstwo tę grupę dziedziczy, więc sygnał wysłany do
ujemnego identyfikatora grupy kończy całe drzewo naraz. Pole pid to zapamiętany identyfikator grupy
procesów, równy identyfikatorowi procesu okna przez wywołanie Setpgid, utrwalony w chwili przejęcia,
gdy jest znany i dodatni. Ubicie posługuje się tym polem, a nie identyfikatorem z os.Process, ponieważ
zwolnienie zeruje go na wartość minus jeden, co dałoby zarazem wyścig danych i policzenie sygnału do
procesu init albo rozgłoszenie do wszystkich procesów systemu. Pole zapisuje się raz, przed jakąkolwiek
współbieżnością, i tylko czyta się je później, więc zwolnienie go nie dotyka.

Metoda ubij, gdy grupa procesu już nie istnieje, bo proces zdążył ją zmienić albo grupy nigdy nie było,
kieruje sygnał wprost do samego procesu, żeby nie zostawić go przy życiu; kod błędu oznaczający brak
procesu mówi, że nie ma już czego ubijać. Metoda posługuje się zapamiętanym identyfikatorem grupy, nie
identyfikatorem z os.Process, z tego samego powodu co przy przejęciu. Wartość identyfikatora nie większa
niż jeden oznacza brak prawidłowego procesu do ubicia i wtedy nie idzie żaden sygnał; straż na wartości
większej niż jeden pilnuje zarazem sygnału do grupy i sygnału bezpośredniego, żeby żaden nie wyrodził
się w sygnał do procesu init albo w rozgłoszenie do wszystkich procesów systemu.
## budowa/server/internal/dane/przekazanie_okna_wiez.go
Więź koordynator-wykonawca nie ma własnej tabeli: mieszka w kolumnie tabeli okien wskazującej okno koordynatora, do której kontrakt odwołuje się wprost. Ten plik dokłada do repozytorium przekazań drogę zapisu i odczytu tej kolumny z poziomu identyfikatora zewnętrznego okna, podczas gdy repozytorium okien czyta i pisze tę kolumnę wyłącznie jako część pełnego wiersza okna po identyfikatorze wewnętrznym, a widok zarządzania oknami operuje na identyfikatorach zewnętrznych pojedynczej więzi, nie całego okna. Cicha zgoda na więź z oknem, którego nie ma, dałaby potwierdzenie relacji, która w rzeczywistości nie powstała, dlatego obie strony więzi są rozwiązywane na identyfikatory wewnętrzne przed zapisem, a nie podzapytaniem w poleceniu aktualizacji, które ciche niedopasowanie zamieniłoby w wartość pustą.

## budowa/server/internal/dane/konta.go
Wzorcem jest rejestr kanałów modelu. Odczyt katalogu nie wynosi odwołania do
poświadczenia. Zapytania tego pliku zwracają wyłącznie znacznik obecności
poświadczenia; samo odwołanie ma jedną, jawnie nazwaną drogę wyjścia w pliku
zapisu kont, a sekretu nie ma w bazie w ogóle.

Katalog rozstrzyga, które konta wolno wziąć do rotacji i w jakiej kolejności.
Kiedy przejść na następne konto, rozstrzyga pula rotacji warstwy iniekcji —
repozytorium nie powiela tamtej logiki.

## budowa/server/internal/dane/aod_wyciszenie.go

Wyciszenie ma własny wiersz, a przypięcia obserwacji nie mają: przypięcie wskazuje proces
telemetrii, który ginie razem z rdzeniem, więc taki wiersz przeżyłby byt, na który wskazuje.
Wyciszenie natomiast wskazuje moduł, kartę sesji albo klasę zdarzeń — byty, które restart
rdzenia przeżywają — i ma sięgać wszystkich powłok Operatora; bez wiersza Operator wyciszałby
w jednej powłoce, a w drugiej sugestie wchodziłyby dalej.

Wyciszenie przeterminowane jest usuwane przy odczycie, nie zegarem w tle: wykaz czyta się
przed każdym ujawnieniem sugestii, a proces budzony cyklicznie po to, żeby zwykle nie zrobić
nic, jest kosztem bez skutku.

Sygnał wyciszony nadal się odkłada: wyciszenie wstrzymuje wyłącznie ujawnienie, nie zapis.
Sito wyciszeń stoi po stronie rdzenia, nie w tym zapytaniu — tu leży wyłącznie zapis i odczyt.

## budowa/server/internal/narzedzia/rozdzielnia.go
Rozdzielnia nie zna ani protokołu MCP, ani gniazda rdzenia: zna wykaz
kontraktu, odwzorowanie komend narzędzi i port do rdzenia, więc daje się
sprawdzić bez procesu po drugiej stronie. Odmowa, komenda nieznana rdzeniowi,
zerwane gniazdo — każda z tych rzeczy wraca do modelu jako czytelny opis
błędu narzędzia, model czyta, poprawia i próbuje dalej.

Rola przekazana do NowaRozdzielnia jest parametrem wymaganym, nie doklejką
z wartością domyślną: zasięg rozstrzyga o tym, co model może zrobić, więc
każdy, kto rozdzielnię składa, ma powiedzieć wprost, w czyim imieniu ona
pracuje. Zasięgiem zwykłym jest zasięg okna.

ZEkspertem stoi osobno od konstruktora, a nie jako kolejny jego parametr:
zasięg eksperta jest jedynym, który potrzebuje wartości, a dokładanie jej
wszystkim pozostałym kazałoby im podawać pustkę bez znaczenia. Dziennik jest
tam wymagany, a nie opcjonalny, bo cena zestawu i każdy brak zawężenia mają
dokądś dojechać; dziennik pusty ucisza je, więc dobór bez dziennika byłby
dokładnie tą cichą degradacją, przeciw której powstał.

Bez zawężenia dołożenia sesji nie mają czego dołożyć: okno bez eksperta ma
pełny wykaz, więc każde dołożenie już w nim stoi. Rozstrzygnięcie, czy
argument dołożeń w ogóle wysłać, należy do strony rdzenia, która wie o oknie
więcej.

Kontekst w Narzedzia wchodzi parametrem, bo w zasięgu eksperta ta droga pyta
rdzeń, a pytanie bez kontekstu nie dałoby się przerwać razem z resztą
procesu; zasięgi pozostałe kontekstu nie tykają.

Narzędzie dołożone przez rolę nie dostaje uzupełnienia okna, i to jest jego
istota, nie przeoczenie: rozszerzenie okna asystenta służy nastawianiu okna
docelowego, a podstawienie okna serwera w brakujące pole okna kazałoby
asystentowi przestawić kanał modelu samemu sobie — dokładnie to, czego zakaz
trzyma te komendy poza wykazem kontraktu. Okno docelowe model wskazuje
jawnie, albo rdzeń odmawia.
## budowa/server/internal/dane/queue.go
Obszar kolejek ma jedno repozytorium rozłożone na dwa pliki wyłącznie dla objętości: ten plik dokłada wykaz kolejek i powiązania obok komend zakładania i akcji obsłużonych w pliku sąsiednim, na tym samym typie repozytorium, nad tą samą bazą i z tym samym dziennikiem akcji. Zawężenia wykazu po sesji i po oknie nie ma celowo: sesja kolejki bywa znana wyłącznie z pamięci powiązań rdzenia, a okna obsługiwane przez kolejkę leżą w tabeli powiązań albo w tej samej pamięci, więc sito po obu tych bytach składa się w rdzeniu, na kolejce kontraktu — sito w zapytaniu SQL milczałoby o kolejkach, które warunek spełniają. Powiązania kolejki dokładają się, nie zastępują: kontrakt zna wyłącznie wiązanie, komendy rozwiązującej nie ma, więc zapis, który cicho zdejmowałby wcześniejsze powiązania, robiłby czynność, o którą nikt nie prosił. Powtórzone powiązanie jest tym samym faktem, a nie drugim, więc kolizja z więzem jednoznaczności nie jest tu błędem, tylko brakiem zmiany.

## budowa/server/internal/session/ubicie_windows.go

Proces okna wraz z całym potomstwem należy do jednego zadania systemowego, więc jedno wywołanie jądra
kończy całe drzewo naraz, zamiast zewnętrznego narzędzia kończącego procesy — bez zależności od narzędzi
zewnętrznych systemu. Uchwyt zadania trzyma zamek: metoda ubij, wywoływana przez obserwatora albo przy
zamknięciu okna, czyta go wtedy, gdy metoda zwolnij, wywoływana przez obserwatora zakończenia procesu,
może go właśnie oddawać; bez zamka byłby to wyścig o pole zadanie. Pole pid to identyfikator procesu
okna utrwalony w chwili przejęcia. Dogląd posługuje się tym polem, a nie identyfikatorem z os.Process,
ponieważ zwolnienie zeruje go na wartość minus jeden, więc czytanie go z gorutyny doglądu byłoby
wyścigiem danych z ubiciem idącym równolegle. Pole zapisuje się raz, w metodzie przejmij, zanim
struktura wyjdzie poza jedną gorutynę, i tylko czyta się je później.

Metoda czekaj otwiera własny uchwyt synchronizujący, żeby nie ruszać uchwytu trzymanego przez
os.Process, który oddaje metoda zwolnij; nieudane otwarcie oznacza, że proces już zniknął. Metoda
posługuje się zapamiętanym identyfikatorem procesu, a nie identyfikatorem z os.Process, z tego samego
powodu co dogląd: zwolnienie zeruje go na wartość minus jeden.

## budowa/server/internal/dane/konta_rotacja.go
Ten plik daje katalogowi dwie rzeczy: uporządkowaną listę kont, które w ogóle
wolno wziąć, oraz trwały ślad tego, co pula już rozpoznała, wraz z miejscem
na decyzję operatora o zawieszeniu konta. Repozytorium nie wygasza
wyczerpania po czasie i nie wybiera konta bieżącego — robi to pula, która
jedyna zna chwilę wywołania.

Konta wyczerpane zostają na liście — o ich pominięciu rozstrzyga pula, która
zna chwilę wywołania i chwilę odnowienia limitu.

## budowa/server/internal/dane/konta_zapis.go
Odwołanie do poświadczenia ma w tym pliku dokładnie dwie drogi — ustawienie
poświadczenia jako wejście i jego odwołanie jako wyjście dla warstwy, która
musi je rozwiązać w magazynie sekretów. Żaden odczyt wykazu ani żadna
odpowiedź kontraktu tą kolumną nie jedzie.

Kolejność niepodana przy zakładaniu konta zostaje nadana jako następna
w obrębie rodzaju — pula rotacji dostaje porządek bez pytania operatora
o liczbę.

## budowa/server/internal/dane/konfiguracja.go
Rozstrzyganie dziewięciu poziomów zasięgu należy do warstwy konfiguracji,
właściciela pojęcia poziomu zasięgu. Brak wiersza oznacza wartość domyślną,
nie odmowę działania, dlatego odczyt zwraca informację „nie ustawiono",
a nie błąd.

Oś jest prostopadła do poziomu: poziom mówi, jak wąsko obowiązuje wartość,
oś mówi, dla czego — dla platformy, dla modelu albo dla konta. Oś pusta
znaczy oś platformy, tak samo jak domyślna wartość kolumny.

Metody kontraktu opisują oś platformy. Oś modelu i konta obsługuje osobne
rozszerzenie kontraktu; obie postaci wypełnia jedna implementacja, więc
drugiego rozstrzygania nie ma.

## budowa/server/internal/narzedzia/rozdzielnia_ekspert.go
Ta droga stoi osobno od rozdzielni zwykłej, bo droga zwykła nie ma prawa się
o nią potknąć: bez trybu eksperta cała maszyneria nie rusza ani razu
i rozdzielnia zachowuje się tak, jak bez niej.

Definicję czyta się przy pierwszym wykazie narzędzi, bo gniazdo do rdzenia
jest leniwe z zamysłem. Odczyt udany zapamiętuje się na czas życia procesu:
ekspert okna nie zmienia się w trakcie tury, a pytanie rdzenia przy każdym
wykazie byłoby ruchem bez treści. Odczyt nieudany nie zapamiętuje się nigdy —
rdzeń bywa niegotowy w chwili startu procesu modelu, więc każdy następny
odczyt wykazu próbuje od nowa; zapamiętana porażka zamieniłaby jedno nieudane
połączenie w oknie bez doboru narzędzi na całą turę.

Zawężenie obowiązuje także przy wywołaniu narzędzia, nie tylko przy złożeniu
wykazu: wykaz zawężony, z którego nadal da się wywołać wszystko, byłby
zawężeniem pozornym, oszczędzającym żetony i nie zmieniającym niczego więcej.
Odmowa dla nazwy spoza podzbioru jest treścią dla modelu, nie usterką procesu.

## budowa/server/internal/dane/aplikacje_warsztat.go

Kontrakt żądania aktualizacji warsztatu nie niesie odpowiednika zakładania migawki znanego
z modułu Developer: Operator nie zakłada tu migawki, tylko nadpisuje stan bieżący pliku
warstwy, dlatego tabela ma jeden wiersz na trójkę okna, warstwy i ścieżki.

Odczyt po kluczu naturalnym wywołuje go komenda aktualizacji warsztatu, która przed
nadpisaniem sprawdza, czy plik już istniał — od tego zależy rodzaj zmiany zgłoszonej
w zdarzeniu zmiany warsztatu (założenie czy aktualizacja).
## budowa/server/internal/dane/role.go
Rola okna mieszka w kolumnie tabeli okien, a więź koordynator-wykonawca w kolumnie sąsiedniej, tam gdzie kontrakt widzi rolę okna i identyfikator okna koordynatora. Ten plik nie zakłada żadnego bytu: dokłada wyłącznie zapis i odczyt tych dwóch kolumn po identyfikatorze, którym posługuje się kontrakt. Osobno od repozytorium okien, bo tamten plik czyta i pisze kolumnę roli wyłącznie jako część pełnego wiersza okna po identyfikatorze wewnętrznym, więc zapis samej roli musiałby wpierw wczytać całe okno wraz z katalogami roboczymi i zapisać je z powrotem, nadpisując po drodze pola, o które komenda przypisania roli nie prosi. Wcielenie roli tu nie ma: nie jest kolumną tego wiersza, tylko wpisem w tabeli ustawień na poziomie zasięgu okna, bo tam zapisuje je warstwa kliencka, a rdzeń sięga po ten sam adres przez repozytorium konfiguracji, zamiast zakładać drugą prawdę o wcieleniu. Metoda dostępu do repozytorium ról, a nie pole struktury, bo rola okna nie jest osobnym obszarem danych, tylko widokiem na dwie kolumny obszaru okien, który zestaw już niesie własnym polem; pole dołożone obok tamtego zapowiadałoby drugie repozytorium tego samego bytu.

## budowa/server/internal/narzedzia/schemat.go
Przekład schematu jest płytki z zamysłem: kontrakt niesie już gotowy typ
schematu, typ elementu tablicy i komplet wartości wyliczenia, wyliczone raz
przez generator przy budowie kontraktu; powtórzenie tamtego rozstrzygania
tutaj byłoby drugim odwzorowaniem tych samych typów.

## budowa/server/internal/dane/aplikacje_wdrozenia.go

Tabela wdrożeń niesie wyłącznie środowisko, strategię, wersję, notatki, adres i odnośnik do
logu zlecenia oraz jego stan — nie prowadzi prawdziwego przebiegu wdrożenia, ponieważ rdzeń
niczego sam nie wdraża, tylko zapisuje to, co dostał od zlecenia.

Zawężenie wykazu wdrożeń do środowiska idzie pustym łańcuchem jako brak zawężenia w jednym
zapytaniu, zamiast dwóch wariantów tekstu SQL sklejanych warunkowo, które rozjeżdżałyby się
przy pierwszej zmianie kolumn.

Liczba pozycji spełniających te same warunki co strona wykazu, ale bez ograniczenia LIMIT,
opisuje rozmiar całej historii wdrożeń, nie rozmiar zwróconej strony.

## budowa/server/internal/dane/komponenty.go
Wiersz komponentu jest kaflem Strefy 2, który wskazuje byt magazynu
modułowego kolumną bytu docelowego. Nie powiela bytu modułowego i nie jest
jego drugą prawdą: kroki automatyki, umiejętności eksperta i pamięć projektu
zostają w swoich tabelach, a to repozytorium nie tyka żadnej z nich.

Konstruktor bierze współdzieloną pamięć zapytań zestawu, tak jak pozostałe
repozytoria pakietu — zamknięcie zestawu zwalnia wyłącznie tę jedną pamięć
poleceń. Połączenia z bazą konstruktor nie bierze i brać nie musi: żaden
zapis tego rejestru nie obejmuje drugiej tabeli, więc transakcji
wielotabelowej tu nie ma.

Kontrakt oznacza pole dołączenia komponentów niczynnych jako niewymagane,
a Strefa 2 pokazuje domyślnie kafle czynne.

Baza nie wstawia własnego znacznika czasu, bo kolumna niesie wartość
kontraktu bez przekładu; dwa zegary dla jednego pola byłyby dwiema prawdami.

Drugi wynik przypisania komponentu mówi, czy przypisanie coś zmieniło —
powtórzenie tego samego przypisania nie dochodzi do skutku i odpowiedź
komendy oddaje wtedy brak zmiany zamiast udawać czynność.

Porządek wykazu komponentów biegnie indeksem tabeli po rodzaju, nazwie
i identyfikatorze, więc kolejność wyświetlania kontraktu jest stała między
wywołaniami.

Warstwa wyższa sprawdza wartość kontraktu, to repozytorium sprawdza skutek
zapisu.

## budowa/server/internal/narzedzia/schemat_roli.go
Kontrakt wylicza parametr narzędzia gotowy — typ schematu, typ elementu
tablicy, komplet wartości wyliczenia — wyłącznie dla komend stojących
w sekcji narzędzi generowanej z kontraktu. Komenda dołożona przez rolę okna
z tej sekcji nie pochodzi, więc taki parametr nie istnieje i nie ma go skąd
wziąć; źródłem staje się wygenerowana struktura żądania, nie ręczny opis, tak
że zmiana pola w kontrakcie zmienia schemat bez dotykania tego pliku.

Ta droga nie daje opisu pola, bo generator Go kładzie opis w komentarzu,
którego odbicie nie widzi, więc pole idzie z opisem pustym. Nie daje też
wykazu wartości wyliczenia: typ nazwany jest w Go zwykłym napisem, a rejestru
wartości pakiet współdzielony nie wystawia, więc pole zostaje napisem, a model
dostaje odmowę rdzenia przy wartości spoza wykazu i poprawia. Nie rozwija też
struktury zagnieżdżonej: pole rodzaju obiektu idzie jako obiekt bez
właściwości, bo rozwijanie w głąb urosłoby do schematu większego niż całe
okno kontekstu, a granicy głębokości kontrakt nie stanowi.

## budowa/server/internal/store/baza.go

Nadmiar równoległych połączeń zamienia rywalizację o zapis w błąd zablokowanej bazy zamiast czekać na
limit czasu zajętości. Skromny limit połączeń trzyma pulę w ryzach, a bezczynne połączenia utrzymuje
ciepłe, żeby dziennik zapisu wyprzedzającego nie był otwierany i zamykany bez końca.
## budowa/server/internal/dane/roundtable.go
Repozytorium nie rozmawia z modelem: wywołanie kanałów uczestników prowadzi rdzeń przez rejestr kanałów, a tutaj leży wyłącznie to, co po debacie zostaje, czyli kto brał w niej udział, o co pytano i co odpowiedziano. Obszar Roundtable urósł ponad jeden plik, więc kontrakt składa się z części: rdzeń debaty stoi w tym pliku, a zdolności dobudowane późniejszymi migracjami leżą we własnych plikach i wchodzą do kontraktu przez zanurzenie, dając jedno repozytorium, jeden kontrakt i tyle plików, ile odpowiedzialności. Rozpoznanie zapisu, który nic nie zmienił, jest osobne od rozpoznania w warstwie sesji, bo tamto miejsce opisuje wiersz numerem klucza głównego i zwraca błąd opisowy, a obszar Roundtable rozpoznaje brak bytu przez porównanie z błędem braku wiersza, dzięki czemu rdzeń oddaje wtedy kod nieznalezienia zamiast usterki wewnętrznej.

## budowa/server/internal/dane/konto_wlasciciela.go
Konto jest jedno i pilnuje tego schemat warunkiem równości identyfikatora
jedynce. Repozytorium nie powtarza tej reguły w kodzie: drugi zapis odbija
się o bazę, a nie o sprawdzenie, które ktoś kiedyś usunie.

Hasła tu nie ma. Tożsamością konta są login i adres e-mail; skrót hasła
leży w sejfie poświadczeń, a wiersz metody uwierzytelnienia niesie do niego
odwołanie.

W tabeli dróg potwierdzenia leży skrót drogi, nigdy sama droga. Kopia bazy
nie daje więc możliwości potwierdzenia cudzej tożsamości — ze skrótu nie
odtworzy się materiału, który poszedł listem.

Rejestracja jest wykonalna raz, więc konto zostawione po nieudanym nadaniu
listu byłoby platformą nie do otwarcia: wejść nie ma czym, bo adresu nikt
nie potwierdził, a założyć drugi raz nie wolno.

Droga nieznana i droga wygasła to dwa różne stany, rozstrzygane przez
warstwę wyżej. Bez sprawdzenia skutku zamknięcia dwa równoległe żądania
z tą samą drogą oba uznałyby ją za ważną.

## budowa/server/internal/dane/asystent.go

Czas w tabeli zlecenia jest liczbą, nie napisem: kolumny czasu niosą milisekundy epoki wprost
jako liczbę całkowitą, bez przekładu przez funkcję formatującą, w odróżnieniu od modułów
trzymających czas tekstem.

ProfilKod nie jest ozdobą wiersza: wykonawca zlecenia biegnie procesem długo po tym, jak
komenda głosowa już odpowiedziała, więc czyta warunki tury z bazy, nie z żądania — profil
trzymany tylko w pamięci procesu nie przetrwałby ani odpowiedzi na komendę, ani restartu
rdzenia w trakcie zlecenia.

Priorytet zlecenia i stan zlecenia to dwie różne czynności Operatora, sterowane osobnymi
metodami, bo zmiana priorytetu nie pociąga za sobą zmiany stanu.

ZakonczZlecenie ustawia stan końcowy i wynik jednym zapisem, żeby okno monitorujące zlecenia
nie zobaczyło przez chwilę stanu bez pasującego wyniku ani wyniku bez stanu — oba pola
pochodzą z jednego zdarzenia, zakończenia tury modelu.

## budowa/server/internal/narzedzia/wpiecie.go
Okno wchodzi argumentem uruchomienia, nie zmienną środowiska, jak już stosuje
most mcp-danaco-pulpit-console: to, co rozstrzyga zasięg jednego wpisu, idzie
argumentem, a to, co dzielą wszystkie wpisy, idzie środowiskiem. Cztery
powody: wpis powstaje osobno dla każdego okna, a proces modelu jest jeden na
okno, więc argument należy do wpisu i różni się wpis po wpisie, podczas gdy
zmienna środowiska należy do procesu i tej rozdzielczości nie ma; zmienną
środowiska dziedziczy każdy proces potomny modelu, a argument nie wychodzi
poza to jedno uruchomienie, i identyfikator okna jako uchwyt do sterowania
platformą nie ma powodu wędrować dalej; argument widać w samej konfiguracji
MCP, więc czytając ten plik wiadomo, którego okna dotyczy, a wpis bez
argumentu byłby dla wszystkich okien identyczny; argument pominięty widać od
razu w wierszu uruchomienia, podczas gdy zmienna pusta jest nieodróżnialna od
nieustawionej. Adres rdzenia idzie drogą przeciwną, środowiskiem, bo jest
wspólny dla całej instalacji i czyta go ten sam pakiet konfiguracji co
w rdzeniu.

Okno puste w Wpis daje fałsz: serwer narzędzi bez okna nie miałby zasięgu,
więc zamiast wpisu bez zasięgu lepiej wpisu nie dokładać wcale, a rozmowa
toczy się wtedy bez sterowania platformą, nie z narzędziami mierzącymi
w nikąd. Powód odmowy jest zwracany, bo ścieżkę binarium oddaje się dopiero
po sprawdzeniu, że plik istnieje: gdy pakiet serwera złożono bez serwera
narzędzi, wpis wskazywałby plik, którego nie ma, a narzędzia sterowania
platformą nie działałyby bez śladu; odmowa musi być powiedziana wprost, bo
brak wpisu i wpis martwy wyglądają dla operatora tak samo, a naprawa jest
inna.

sciezkaProgramu szuka binarium obok binarium, które właśnie pracuje, bo
rdzeń i serwer narzędzi wychodzą z jednego budowania i jadą w jednym pakiecie
instalacyjnym, więc stoją w tym samym katalogu. Gdy obok go nie ma, albo gdy
miejsca bieżącego procesu nie da się ustalić, zostaje ścieżka wyszukiwania
systemu, droga uruchomienia z procesu stojącego w katalogu tymczasowym.
Dopiero gdy zawiodą obie drogi, funkcja zwraca błąd zamiast ścieżki: ścieżka
zmyślona jest gorsza od braku wpisu, bo o braku wpisu da się powiedzieć.

brakBinarium jest osobnym typem, a nie samym napisem, żeby wywołujący mógł
ten jeden przypadek odróżnić od pozostałych odmów: to jedyna odmowa, która
znaczy, że produkt zbudowano lub spakowano niekompletnie, a nie że dane okno
nie ma zasięgu. Droga wskazana w komunikacie błędu jest jedyną, którą serwer
narzędzi w produkcie powstaje: instalacja Operatora nie niesie ani rdzenia,
ani serwera narzędzi, więc oba stoją wyłącznie w pakiecie serwera.

## budowa/server/internal/dane/kopiowanie_sesji.go
Kopiowanie idzie po wierszach bazy, nie po rejestrze żywym: kopia ma być
wiernym odbiciem zapisu, a rejestr żywy niesie wyłącznie sesje otwarte.
Sesja sprzed restartu ma dać się skopiować tak samo jak ta z bieżącej pracy.

Odwzorowanie okien źródłowych na docelowe podaje wywołujący — to on nadaje
nowym oknom identyfikatory rdzenia i zna ich wiersze. Warstwa danych nie
zakłada okien sama, bo cykl życia okna należy do pakietu sesji.

Wiadomość kopiowana traci powiązanie z oknem źródłowym: w kopii nie ma bytu,
na który mogłoby wskazywać, a wskazanie na okno oryginału byłoby więzią
między dwiema niezależnymi sesjami. Funkcja zwraca liczbę skopiowanych
wiadomości — jedyną miarę, po którą sięga wywołujący; struktura wyniku
byłaby typem bez odbiorcy.

## budowa/server/internal/dane/asystent_dziennik.go

Notatka pusta przy nadawaniu wyróżnienia zostawia powód zastany: wyróżnienie ponowione bez
własnego słowa nie ma prawa skasować zdania zapisanego wcześniej, w odróżnieniu od zdjęcia
wyróżnienia, które kasuje powód razem ze znacznikiem.
## budowa/server/internal/dane/roundtable_decyzja.go
Zmiana wagi jednego kryterium przestawia wynik każdego wariantu naraz, więc kolumna z wynikiem ważonym rozjechałaby się z ocenami przy pierwszym pominięciu przeliczenia; z tego powodu wynik ważony nie jest przechowywany, tylko liczony przez rdzeń przy odczycie.

## budowa/server/internal/store/katalog_akcji_test.go

Panel akcji i siatka szybkich akcji nie mają własnej listy pozycji — biorą ją komendą listowania z
tabeli akcja. Wiersz katalogu wskazuje komendę kolumną komenda, a kolumna ta nie jest kluczem obcym
i być nim nie może, ponieważ kontrakt mieszka w osobnym pliku, nie w bazie; baza przyjmie więc każdą
nazwę, także nazwę komendy, której nie ma. Skutek takiego wiersza u operatora: kontrolka jest, daje
się nacisnąć i wraca odmową nieznanej komendy. Jest to atrapa — element, który obiecuje czynność, a nie
ma za sobą ani jednego wykonawcy, co jest kłamstwem odwrotnym do wykazu braków, który już raz mówił
o brakach, których nie było. Dlatego kontrola stoi w tym pliku, a nie w kodzie rdzenia: pyta o stan po
pełnym przejeździe migracji, więc obejmuje każdy wiersz katalogu, także wiersz wniesiony migracją, która
jeszcze nie istnieje. Sprawdzian wypada niepomyślnie także wtedy, gdy katalog zostanie opróżniony, bo
zaczyn akcji jest treścią produktu, nie danymi przykładowymi.

Zapory nie wolno osłabić wykazem wyjątków: pozycja bez pokrycia w kontrakcie nie ma stanu przejściowego
„jeszcze nie", dopóki komendy nie ma, kontrolki też nie ma być, a wiersz dochodzi migracją razem
z komendą. Funkcja pozycjeBezPokrycia liczy do tego samego wykazu również pozycję bez ani jednej
komendy, ponieważ kontrolka bez komendy i kontrolka z komendą nieistniejącą kończą się u operatora tym
samym, czyli naciśnięciem bez skutku; wyliczenie stoi osobno od sprawdzianu po to, żeby dało się je
nakarmić wierszem, którego w katalogu nie ma.

Sprawdzian, który przechodzi na katalogu zdrowym, ale przeszedłby też na katalogu z atrapą, jest
sprawdzianem pozornym, a pozorny sprawdzian jest gorszy od jego braku, bo świeci zielono i nikt nie
patrzy dalej. Droga niepomyślna jest tu mierzona wprost: wiersz wskazujący komendę, której nie ma, oraz
wiersz bez komendy muszą wyjść z wyliczenia oba.

Test ikon pilnuje drugiej połowy tej samej obietnicy: pozycja z ikoną, której nie ma w zestawie klienta,
wychodzi w oknie kontrolką bez znaku, pustym prostokątem nieczytelnym dla operatora. Zestaw czyta się
z plików źródłowych ikon, ponieważ są jedyną prawdą o tym, co klient umie narysować. Sprawdzian stoi
i czeka, dopóki klient nie ma zestawu ikon, ponieważ zestaw wchodzi do klienta wraz z ramą aplikacji,
a przed nim nie ma czego czytać; zastępnika nie ma i być nie może, bo zestaw poprzedniego klienta jest
materiałem do przeszczepu, nie tym, co obecny klient umie narysować, a katalog ikon wkompilowany
w rdzeń jest materiałem komend na innej siatce nazw, nie zestawem kontrolek okna — miara wzięta
z któregokolwiek z nich świeciłaby zielono, nie mierząc okna. Warunkiem powrotu jest katalog wskazany
stałą sciezkaZrodelIkon; pominięcie jest warunkowe, więc sprawdzian wraca sam w chwili, gdy zestaw
stanie.

Warunek powrotu w funkcji zrodlaIkonKlientaStoja stoi osobno od funkcji nazwyIkonKlienta, bo dwa stany
trzeba odróżnić: zestawu jeszcze nie ma, co jest pominięciem, oraz zestaw jest, ale nie daje nazw, co
jest niepowodzeniem — zmienił kształt albo sprawdzian czyta niewłaściwy katalog. Odczyt nazw idzie po
kluczach zapisu obiektu, ponieważ nazwa ikony stoi w tych plikach jako klucz wcięty dwoma znakami
odstępu; rozbiór składni języka źródłowego byłby tu kodem, który sam może się mylić, a wzorzec klucza
wystarcza, bo pliki źródeł mają jeden kształt pilnowany przez formater klienta.

## budowa/server/internal/narzedzia/wpiecie_test.go
Odmowa wskazująca plik, którego nie ma, jest gorsza od odmowy milczącej: wysyła
Operatora po nic. Sprawdziany tego pliku nie poprzestają na czytaniu napisu
odmowy — biorą z niego ścieżki i mierzą, czy te ścieżki leżą w drzewie budowy,
bo sam napis mógłby wskazywać cokolwiek i nadal wyglądać poprawnie.

Bez potwierdzenia miejsca w korzenBudowy sprawdzian mierzący obecność pliku
nie odróżniłby pliku nieobecnego od pomiaru wykonanego w złym miejscu: obie
drogi kończą się tym samym błędem odczytu pliku, a znaczą co innego.

Ścieżkę wyszukiwania systemu w powodBrakuBinarium podmienia się na katalog
pusty, żeby wyszukanie binarium nie miało czego znaleźć; obok binarium
sprawdzianu, które stoi w katalogu tymczasowym budowania, serwera narzędzi
też nie ma. Gdyby mimo to Wpis oddał ścieżkę, sprawdzian pada zamiast przejść:
odmowy, której nie było, nie wolno uznać za odmowę zbadaną.

## budowa/server/internal/narzedzia/wykaz.go
Ani jedna nazwa narzędzia, ani jeden opis, ani jedno pole schematu nie są
zapisane w tym pakiecie: wszystko czyta się z funkcjami i odwzorowaniem
wytworzonymi z kontraktu. Dopisanie komendy do sekcji narzędzi kontraktu
powiększa ten serwer bez zmiany choćby jednej linii kodu, i tak samo działa
w drugą stronę: wykreślenie komendy odbiera modelowi narzędzie. Drugiego
wykazu nie ma z zamysłem — wykaz własny rozjechałby się z kontraktem, gdy
tylko kontrakt urośnie.

Grupa narzędzia jest polem danych, bo te same grupy są potem gałęziami
drzewa wyboru u Operatora i jednostką doboru narzędzi eksperta.

Zasięg okna roboczego w WykazZasiegu oddaje sam wykaz kontraktu, bez
dokładania pozycji roli.

## budowa/server/internal/store/macierz_terminal_test.go

Macierz dostępności modułów stawia moduł Terminal jako widoczny w CodeStudio, z własnymi oknami kart
terminala, konsoli wyjścia i monitora procesów. Migracja 080 wygasza go ustawieniem niewidoczności
z uzasadnieniem, że Terminal nie jest samodzielnym modułem — rozstrzygnięciem przeciw jego dostawie.
Sprawdzian zapisuje ten rozjazd tak, żeby nie zniknął po cichu ani nie pogłębił się po cichu: wypada
niepomyślnie, dopóki Terminal jest ukryty w CodeStudio, a gdy zostanie odsłonięty zgodnie z macierzą,
każe się odwrócić w straż. Naprawa należy do warstwy modułów, do macierzy widoczności środowisko-moduł,
nie do tego pliku — sprawdzian pomiar utrwala, nie rozstrzyga.
## budowa/server/internal/dane/roundtable_glosowanie.go
Wyniku agregacji tu nie ma i być nie może: liczy go rdzeń z głosów przy każdym odczycie, bo głos może dojść po pierwszym wyliczeniu. Repozytorium oddaje materiał, czyli głosowanie, warianty i głosy, a nie wnioski wyciągnięte z niego. Powtórne oddanie głosu zastępuje poprzedni: zmiana zdania w otwartym głosowaniu jest czynnością dozwoloną, a dwa głosy tej samej osoby nie są. Otwarcie głosowania zakłada je wraz z wariantami w jednej transakcji, bo głosowanie bez wariantów byłoby pytaniem bez odpowiedzi do wyboru.

## budowa/server/internal/dane/library_kolekcje.go
Plik biblioteki leży w osobnym pliku, wersje w kolejnym osobnym pliku.
Kolekcja nieznana to co innego niż plik nieznany: nie ma dokąd przypisywać,
więc wywołanie wraca jako błąd braku wiersza i adapter odmawia wprost. Pusta
lista dawałaby przypisaniu kolekcji powodzenie z zerem przypisań, czyli
potwierdzenie czynności, która się nie odbyła.

Ustawianie etykiet nadsyła komplet etykiet pliku, nie różnicę — tak samo jak
zapis kroków automatyzacji podmienia komplet kroków. Ustawienie etykiet
usuwa więc zastane etykiety i wstawia nadesłane w jednej transakcji.

## budowa/server/internal/dane/library_wersje.go
Wersja jest własnym bytem, nie polem licznika: każdy wiersz niesie własną
treść, sumę kontrolną i autora, dlatego zapis wersji nie nadpisuje niczego,
tylko dokłada wiersz historii.

Plik po zmianie zwraca osobna metoda odczytu pliku, żeby nie duplikować tu
kształtu struktury pliku biblioteki, którego ten plik nie deklaruje.

## budowa/server/internal/dane/library_wersje_zapis.go
Osobny plik od pliku odczytu wersji: tamten plik odpowiada za odczyt
historii i za przywrócenie wersji zastanej; ten za jej dołożenie. Rozdział
idzie wzdłuż odpowiedzialności, a nie wzdłuż tabeli — polecenie SQL obu
stron jest to samo i mieszka nadal w pliku odczytu, żeby nie było dwóch
prawd o jednym poleceniu.

Dołożenie jest nierozdzielne. Wstawienie wiersza historii i przestawienie
pliku macierzystego na tę wersję to jedna zmiana stanu: plik, którego
wskaźnik wersji bieżącej wskazuje wiersz nieistniejący, albo historia
z wersją, której plik nigdy nie przyjął, to schemat rozjechany w połowie.
Stąd transakcja, tak samo jak przy przywracaniu wersji.

Liczba wersji poprzednich jest oddana wołającemu. Tabela nie ma kolumny
numeru wersji — porządek historii daje sortowanie malejące po dacie
utworzenia i identyfikatorze. Rdzeń, który chce nazwać wersję jej
kolejnością, bierze ją z odczytu wersji; ten zapis niczego nie numeruje, bo
numer nie jest tu bytem trwałym.

## budowa/server/internal/narzedzia/zasieg_eksperta.go
Zasięg eksperta to podzbiór wykazu kontraktu wskazany przez definicję eksperta,
a nie wykaz okna; dwa pozostałe zasięgi wynikają z roli okna i tylko dokładają.
Zawężenie jest potrzebne, bo pełny wykaz kontraktu złożony w kształt tools/list
waży rzędu dziesiątek tysięcy żetonów zjadanych, zanim padnie pierwsze słowo
zadania; przy nadmiarze narzędzia docierają do modelu bez opisów, model widzi
same nazwy i po nie nie sięga.

Kod eksperta przychodzi przełącznikiem uruchomieniowym, czytanym raz przy
starcie, z wpisu danaco ułożonego przez rdzeń. Model nie ma czym o zasięg
poprosić ani go poszerzyć: nie ma narzędzia zmieniającego zasięg, a wpis
powstaje zanim proces modelu wystartuje. Rdzeń bierze kod eksperta z tego
samego pola ustawień okna, które nakłada eksperta na okno.

Zasięg eksperta spada zawsze na model roboczy, nigdy na asystenta klawiatury
Operatora: okno eksperta jest ręką modelu wykonującego zlecenie, nie ręką
Operatora.

Kod pusty w ArgumentyEksperta nie daje żadnych argumentów, ani nazwy zasięgu:
zasięg eksperta bez eksperta nie jest zasięgiem węższym, tylko zasięgiem bez
treści; wpis z samą nazwą roli kazałby serwerowi meldować brak przy każdym
odczycie wykazu zamiast po prostu pracować w zasięgu okna.

## budowa/server/internal/narzedzia/zasieg_okna.go
Wpis danaco w konfiguracji MCP powstaje osobno dla każdego okna i niesie jego
identyfikator, więc narzędzie zawsze wie, z którego okna przyszło wywołanie.

Wskazanie jawne okna zostaje: kiedy model podaje pole okna sam, wartość idzie
do rdzenia bez zmiany, bo pętla koordynator-wykonawca polega na tym, że okno
koordynatora wysyła wiadomość do okna wykonawcy, a opisy narzędzi w kontrakcie
mówią to wprost. Podmienianie wskazania jawnego na własne okno zamknęłoby tę
pętlę i rozminęło serwer z kontraktem, który go opisuje. Pole rozpoznaje się
po nazwie z deklaracji kontraktu, nie po własnym wykazie komend okna, więc
narzędzie bez tego pola przechodzi nietknięte.

Wartość wskazania okna innego rodzaju niż napis w brakWskazaniaOkna zostaje
nietknięta, bo jest wskazaniem wadliwym, a orzekanie o kształcie treści
żądania należy do rdzenia, nie do rozdzielni; podmiana takiej wartości na
własne okno ukryłaby pomyłkę modelu.
## budowa/server/internal/dane/roundtable_konsensus.go
Wersja stanowiska jest wpisem, nie licznikiem: licznik w kolumnie wersji tabeli stanowiska mówi, ile redakcji było, ale porównać dwie redakcje da się dopiero wtedy, gdy każda z nich została zapisana osobno. Powtórny zapis tej samej wersji nie jest błędem: stanowisko odczytywane wielokrotnie bez zmiany treści nie podbija licznika, więc wersja bieżąca pozostaje spójna z ostatnim zapisanym wpisem.

## budowa/server/internal/dane/macierz.go
Wiersze wnoszą migracje schematu — repozytorium ich nie zakłada. To osobny
byt obok repozytorium modułów: tamto repozytorium czyta moduły, a macierzy
używa wyłącznie jako filtra. Tu bytem jest sam wiersz macierzy — z parą
kodów, kolejnością i widocznością — którego tamten kształt nie umie oddać.
To nie druga prawda o module: jedna tabela, dwa różne pytania.

Kody wchodzą do wiersza macierzy razem z identyfikatorami, bo czytelnik
macierzy prawie zawsze potrzebuje kodu, a nie numeru wiersza — a drugie
zapytanie po słownik byłoby powrotem do pętli wielokrotnych zapytań, którą
ten byt właśnie znosi.

Zapisu w tym repozytorium nie ma świadomie. Kontrakt platformy nie
definiuje ani jednej komendy zmieniającej macierz, więc metoda zapisu nie
miałaby drogi wywołania — a byt bez drogi wywołania jest atrapą.

Moduł nieobecny w wyniku odczytu widocznych kodów środowisk nie ma okna
modułowego w żadnym środowisku — jest dostępny wyłącznie ze strony głównej.

Jedno zapytanie zamiast zapytania na każde środowisko: macierz jest mała,
ale czyta ją każde wejście na stronę główną, wykaz środowisk, wejście do
środowiska i wykaz modułów, więc wielokrotne zapytania płaciłyby się przy
każdym wejściu operatora.

## budowa/server/internal/dane/auth.go

Wiersz metody uwierzytelnienia nie jest kontem — kont użytkownika platforma nie prowadzi —
ani poświadczeniem kanału modelu, które mieszka w osobnej tabeli konta i z bramką nie ma
związku. Skrót hasła składa rdzeń i on kładzie go w sejfie poświadczeń.

Kotwica to hasło bramki: brak jej wiersza otwiera jednorazową wykonalność rejestracji.

UniewaznijSesjeBramkiPoza: pusty skrót sesji własnej znaczy unieważnienie wszystkich sesji
czynnych.

UrzadzeniaKonta niesie też informację, czy urządzenie ma dziś ważny token.

Urządzenie konta nie ma bytu trwałego w osobnej tabeli i nie musi go mieć: urządzeniem konta
jest to, które kiedykolwiek weszło, a to wiedzą sesje bramki; osobna tabela urządzeń
wymagałaby sprzątania wierszy, których nic już nie dotyczy.

## budowa/server/internal/store/migracje_test.go

Krok migracji wykonany błędnie zostaje w bazie na zawsze, bo kroku nie da się cofnąć. Dlatego mierzone
jest tu nie to, czy przejazd przechodzi, lecz cztery obietnice, na których stoi cała reszta: pełny
przejazd od zera, powtarzalność, niezmienność treści kroku już zastosowanego oraz atomowość kroku
nieudanego. Sterownik jest czystym kodem języka Go, więc każdy sprawdzian zakłada własny plik bazy
w katalogu tymczasowym i przejeżdża komplet migracji od nowa; koszt tego przejazdu jest na tyle mały,
że nie opłaca się dzielić bazy między sprawdzianami, bo baza dzielona zamieniłaby je w jeden sprawdzian
zależny od kolejności.

Ciągłość numeracji kroków migracji nie jest wymagana i nie jest sprawdzana w teście numeracji — luki
powstają przy pracy równoległej i są z zamysłu.

## budowa/server/internal/narzedzia/zasieg_roli.go
Serwer narzędzi należy do jednego okna, a wpis danaco w konfiguracji MCP
powstaje osobno dla każdego okna, więc zasięg narzędzi wiąże się z oknem.
Wykaz narzędzi kontraktu nie niesie komend przestawiania konfiguracji sesji
ani kanału modelu: model nie przestawia sobie własnego wyposażenia, co jest
właściwe dla okna roboczego. Okno asystenta ustawia jednak konfigurację
zlecenia i wybiera kanał modelu — nie na sobie, lecz na oknie docelowym,
w którym ma pracować model wykonujący zlecenie.

Zasięg jest funkcją roli okna, nie prośby modelu. Rozszerzenie wynika z roli
okna zapisanej w rdzeniu i wchodzi do wpisu MCP w chwili jego składania,
zanim proces modelu wystartuje; model nie ma czym o rozszerzenie poprosić,
bo nie ma narzędzia zmieniającego zasięg, a przełącznik zasięgu czyta się raz,
przy uruchomieniu serwera, z wpisu ułożonego przez rdzeń.

Rozszerzenie nie tyka komend zastrzeżonych Operatorowi — punkty dostępu,
nadania, konta, tożsamość — ani warstwy połączenia klienta; te zostają poza
wykazem w każdym zasięgu.

Wykaz rozszerzenia jest polityką, nie kopią danych kontraktu: kontrakt nie zna
pojęcia okna asystenta, sekcja narzędzi jest listą płaską, bez wymiaru roli,
więc nie ma w niej danych, z których dałoby się ten podzbiór wyprowadzić.
Dlatego wykaz jest wymieniony z nazwy, opatrzony powodem i sprawdzany przy
budowie zasięgu: nazwa, która przestanie być komendą kontraktu albo wejdzie
do wykazu narzędzi, znika z rozszerzenia zamiast po cichu wisieć.

Schemat wejścia komend roli powstaje z wygenerowanej struktury żądania, bo dla
komendy spoza sekcji narzędzi kontrakt nie wylicza parametrów; pola mają więc
typy, ale nie mają opisów ani wykazu wartości wyliczeń.

Rozpoznawane są trzy zasięgi w RozpoznajZasieg; trzeci, ekspercki, definiuje
osobny plik zasięgu eksperta. Zasięg węższy jest bezpiecznym domyślnym,
a rozszerzenia nie dostaje się przez pomyłkę w napisie przełącznika.

Warunek istnienia w kontrakcie i stania poza wykazem narzędzi w narzedziaRoli
chroni odpowiednio przed nazwą, która z kontraktu wypadła, i przed podwójną
pozycją tego samego narzędzia, gdyby kontrakt kiedyś wciągnął tę komendę do
wykazu sam; oba warunki liczy ta sama funkcja, którą posługuje się odmowa —
jedno źródło rozstrzygnięcia, co stoi poza wykazem.

## budowa/server/internal/store/nastawy_przesiewu_i_mowy_test.go

Bez wiersza w katalogu ustawień zdolność dalej działa, ponieważ rozstrzyganie nastawy czyta zapis
niezależnie od katalogu definicji, ale zapis ustawienia odmawia klucza spoza katalogu, a okno
konfiguracji wystawia wyłącznie pozycje katalogu. Brak wiersza znaczy więc dokładnie tyle: nastawa jest
ustawialna wyłącznie ręcznym zapisem do bazy. Sprawdzian pilnuje jednego i drugiego naraz — że wiersz
jest i że jego wartość domyślna mówi to samo, co stała pakietu, z którego liczy silnik.

Wartość domyślną odbiorca dostaje z bazy przez rozstrzygacz zasięgu, a stała pakietu mowa wchodzi tam,
gdzie rozstrzygacza nie ma. Rozjazd daje dwie odpowiedzi na pytanie, gdzie leżą wagi, zależne od drogi
wywołania, a rozpoznać go można dopiero po tym, że transkrypcja pobiera drugą kopię wag zamiast
wystartować. Katalog pusty oznacza wtedy, że widoczność wag zależy od tego, na czyim koncie stoi
proces, więc wartość pusta jest tu regresją, a nie wyborem.
## budowa/server/internal/dane/roundtable_sklad.go
Zespół jest kopią składu, nie odwołaniem do niego: skład okna zmienia się po zapisaniu zespołu, bo uczestnicy dochodzą, wypadają i zmieniają rolę, a zespół ma zostać taki, jaki był w chwili zapisu. Odwołanie do wierszy składu zamiast kopii dałoby zespół, który wnosi do nowego okna stan cudzego okna z bieżącej chwili zamiast stanu zapamiętanego. Zmiana tożsamości uczestnika nie dotyka wyciszenia ani kolejności, bo obie te własności są czynnościami moderatora, a zapis tożsamości opisuje wyłącznie samego uczestnika. Usunięcie uczestnika ze składu nie usuwa jego wypowiedzi: zapis tury nie ma prawa się zmienić dlatego, że mówca wypadł ze składu. Zapis zespołu utrwala go razem ze składem w jednej transakcji, bo zespół z połową składu byłby układem, którego nikt nie zapisywał.

## budowa/server/internal/dane/moduly.go
Wiersze wnosi zaczyn schematu — repozytorium ich nie zakłada.

Rodzaj modułu przybiera jedną z czterech wartości: środowisko robocze
(czat, narzędzia, okna pomocnicze), kompozytor (wytwórnia elementów
używanych w innych modułach), repozytorium plików (menedżer zasobów bez AI)
albo sekcja konfiguracyjna. Pole jest wskaźnikiem, nie napisem: brak wartości
znaczy „rodzaju nie ustalono" i jest odróżnialny od każdej z czterech
wartości. Kolumna nie ma warunku sprawdzającego — silnik bazy nie umie go
dołożyć zmianą schematu po fakcie — więc zbioru wartości pilnuje treść
migracji, nie schemat.

Fakt konfigurowalności na stronie głównej jest niezależny od rodzaju
modułu: modułów, które go niosą, nie łączy jeden rodzaj i nie dałoby się
ich z rodzaju wyprowadzić. Na zewnątrz wychodzi jako pole kontraktu
przekładane warstwą nawigacji.

## budowa/server/internal/dane/auth_sesje_bramki.go

Sesja bramki nie jest kartą sesji: karta sesji opisuje pracę — rozmowę, okna, kolejki,
projekt — i istnieje niezależnie od tego, czy ktokolwiek się zalogował, a wygaśnięcie
wejścia karty pracy nie zabiera.

Wykaz urządzeń liczy, czy urządzenie ma dziś ważny token, czyli czy unieważnienie sesji
miałoby co odebrać.

Unieważnienie sesji bramki: pusty napis w argumencie zdejmuje wyłączenie sesji bieżącej,
więc jedno zapytanie obsługuje oba warianty — poza sesją bieżącą i wszystkie — zamiast
rozjeżdżać się na dwie osobne ścieżki kodu.

## budowa/server/internal/narzedzia/zastrzezenia.go
Komendy zastrzeżone Operatorowi obejmują zakładanie i zmianę punktów dostępu,
nadań, kont oraz treści tożsamości; odczyt tych rejestrów model ma, zapisu
nie. Granica jest strukturalna, nie wpisana: rozdzielnia zna wyłącznie
odwzorowanie komend narzędzi, a w nim komend zastrzeżonych po prostu nie ma,
więc żadne wywołanie nie ma jak do nich dojść, choćby model podał nazwę
wprost. Ten plik dokłada do tego drugie: nazwie spoza wykazu odpowiada
czytelna odmowa mówiąca, o którą komendę chodzi, zamiast milczącej odmowy,
po której model próbowałby dalej.

Własnego wykazu zastrzeżeń tu nie ma: zbiór powstaje z różnicy komend
kontraktu i komend narzędzi, więc przesunięcie granicy w kontrakcie przesuwa
ją tutaj w tej samej chwili.

Odmowa w zbierzKomendyPozaWykazem schodzi na wariant ogólny, gdy wzorca nie
da się odczytać, zamiast zgadywać nazwy własnym wykazem.

Wzorzec nazw narzędzi mieszka w kontrakcie, lecz generator Go go nie
wyprowadza. Zapisanie go literałem byłoby drugim źródłem prawdy, które
rozjedzie się przy pierwszej zmianie notacji, dlatego wzorzec jest odczytany,
a nie założony, i potwierdzony na każdej deklaracji wykazu.

## budowa/server/internal/store/nastawy_wiedzy_test.go

Wartość domyślną odbiorca dostaje z bazy: rozstrzygacz zasięgu oddaje wartość domyślną definicji
ustawienia jako rozstrzygnięcie o tym pochodzeniu, a adapter wiedzy nanosi je na komplet nastaw. Stałe
pakietu wiedza wchodzą tam, gdzie rozstrzygacza nie ma. Rozjazd między jednym a drugim nie wywraca
niczego przy starcie — daje dwie różne odpowiedzi na pytanie, czym rdzeń liczy, zależne od drogi
wywołania, a rozpoznać to można dopiero po pustym wyniku wyszukiwania. Sprawdzian wiąże więc obie
prawdy w jedną.

Świeże wdrożenie ma wystartować na wagach rozłożonych na maszynie odbiorcy, więc pusty katalog wag jest
tu regresją, a nie wyborem.

## budowa/server/internal/dane/mobile.go
Warstwa mobilna nie ma własnej tabeli. Proces mobilny jest wpisem rejestru
telemetrii postępu trzymanego w pamięci; sesje, okna, procesy i kolejki mają
swoich właścicieli, a ten plik dokłada wyłącznie odczyt.

Metody siedzą na repozytoriach okien i kolejek, bo tabela ma jednego
właściciela. Osobne repozytorium mobilne z własnymi zapytaniami nad tymi
tabelami byłoby drugim czytelnikiem cudzych tabel i rozjechałoby się
z właścicielem przy pierwszej zmianie słownika stanów. Osobny jest wyłącznie
plik, żeby widać było, po co te rachunki powstały.

Czynność bytu mierzy się tu stanem spoza stanów końcowych: sesja czynna to
stan czynny albo wstrzymany, nigdy zakończony ani archiwalny. Kolejka idzie
tą samą miarą: czynna jest bezczynna, pracująca i wstrzymana, końcowe są
zatrzymana i wyczerpana. Okno komunikacji zna dwa stany, więc liczą się
wiersze w stanie otwartym.
## budowa/server/internal/dane/roundtable_stanowisko.go
Stanowisko powstaje ze złożenia wypowiedzi tury, więc odczyt bez nowej wypowiedzi daje treść tę samą, a podbijanie wersji przy każdym otwarciu okna zamieniłoby licznik wersji w licznik odczytów. Bez warunku broniącego redakcji Operatora pierwsze otwarcie panelu stanowiska po redakcji wracałoby do zapisu tur i kasowało pracę Operatora. Redakcja Operatora podnosi wersję zawsze, także wtedy, gdy treść wyszła ta sama, bo zapisanie tej samej treści jest czynnością zamierzoną, nie powtórzeniem bez skutku.

## budowa/server/internal/podagenci/domkniecie.go
Praca podagenta biegnie w tle, a jej odwołanie leży w wykazie prac pod kodem
podagenta. Kiedy praca się kończy, także gdy kończy błędem, wykaz zwalnia się
automatycznie. Zapis stanu końcowego idzie po tym i sam też potrafi paść,
choćby na chwilowe zajęcie bazy. Zostaje wtedy wiersz w stanie niekońcowym
bez uchwytu do pracy, a zatrzymanie podagenta odwołuje wyłącznie po uchwycie:
uchwytu nie ma, więc melduje brak działania i wiersza nie rusza. Podagent
stoi w stanie oczekującym albo biegnącym do końca życia bazy.

Odpowiedzią jest domknięcie, a nie zapisywanie stanu zatrzymanego zawsze, bo
meldunek o braku działania niesie potrzebną wiedzę: Operator zatrzymujący
wielu podagentów ma wiedzieć, których zdążył zatrzymać, a którzy skończyli
sami. Domknięcie robi rzecz węższą — bierze wyłącznie wiersze niekońcowe,
których pracy nikt nie prowadzi, i przestawia je na stan końcowy. Wiersz
zakończony zostaje nietknięty.

Zapis stanu jest uparty, nie jednokrotny: nieudany wraca jeszcze kilka razy
z rosnącym odstępem, bo rywalizacja o zapis w bazie jest chwilowa, a jedna
nieudana próba zamieniłaby ją w stan trwale nieprawdziwy.

Praca oddana przed urwaniem jest ważniejsza niż wyjaśnienie, dlaczego się
urwała — dlatego wyjaśnienie domknięcia wchodzi wyłącznie do pola wyniku
pustego, warunek stawia funkcja pomocnicza, nie samo zapytanie.

Kontekst zerwany w ZapiszStan kończy ponawianie od razu: zatrzymany rdzeń nie
ma po co czekać na bazę, a zapis pod kontekstem odwołanym i tak nie ma prawa
się udać, dlatego wołający podaje kontekst życia rdzenia, nie kontekst
odwołanej pracy. Trwałość pusta kończy się błędem nazywającym jej brak, a nie
meldunkiem udanego zapisu.

DomknijNieczynnych woła się po stwierdzeniu bezczynności, nie zamiast niego:
wykaz podaje wołający i to on odpowiada za to, że praca tych podagentów
naprawdę nie biegnie — sam przegląd żywotności tego nie mierzy. Błąd jednego
wiersza nie przerywa pozostałych, bo domknięcie części wykazu jest lepsze niż
odmowa domknięcia całości; pierwszy napotkany błąd wraca po przejściu całego
wykazu.

Zapis stanu wynikowego używa reguły pierwszej wartości niepustej, a wartość
niepusta nadpisałaby więc pracę, którą podagent zdążył oddać; skoro zapytanie
warunku nie stawia, stawia go wołający.

## budowa/server/internal/dane/nagrania_mowy.go
Bajty leżą na dysku, w katalogu danych rdzenia; tutaj mieszka wyłącznie
wiersz opisujący jedno nagranie. Rozdział jest zamierzony: odnośnikiem
nagrania w całym produkcie jest ścieżka pliku, bo taką przyjmuje transkrypcja
mowy i taką oddaje synteza mowy panelu tłumaczenia. Drugi rodzaj odnośnika
oznaczałby przekład w każdym miejscu styku.

Rejestr jest też wykazem tego, co rdzeń sam wystawił — a więc granicą
odsłuchu: pobranie nagrania oddaje bajty z tego wykazu, a nie dowolnego
pliku, którego ścieżkę ktoś przyśle.

## budowa/server/internal/store/odwzorowanie_kontraktu_test.go

Odwzorowanie kontraktu na schemat bazy ma być jeden do jednego i żyć wyłącznie w kontrakcie, ale ta
deklaracja jest dziś obietnicą w pliku kontraktu, którą nic nie egzekwuje: kolumna może zostać
przemianowana kolejną migracją, wartość może zostać dodana do kontraktu i pominięta w warunku
sprawdzającym, a obie zmiany przechodzą kompilację po obu stronach i wychodzą dopiero zapisem
odrzuconym przez bazę u odbiorcy. Źródłem prawdy jest tu schemat wygenerowany przejazdem migracji, nie
treść plików migracji: ta sama nazwa kolumny występuje w kilkunastu tabelach naraz, więc odczyt z plików
nie rozstrzygnąłby, o którą tabelę chodzi, a odczyt ze słownika schematu rozstrzyga.

Wykaz odwzorowaniaRozeszlyeSieZeSchematem jest zaporą, nie zgodą — sprawdzian wypada niepomyślnie także
wtedy, gdy rozjazd zostanie usunięty, a wiersz zostanie w wykazie. Wpis ProgressStatus wskazuje na
kolumnę stanu tabeli procesów sesji, która powstała w kroku zakładającym okna, a odeszła w kroku
zdejmującym sieroty transportu, ponieważ rejestr procesów żyje w pamięci rdzenia i tam jest jego jedyne
miejsce; odwzorowanie w kontrakcie zostało po tabeli, której nie ma, a rozstrzygnięcie należy do
kontraktu, nie do sprawdzianu, bo to plik kontraktu jest źródłem prawdy nazw, a nie schemat. Dalsza
pozycja wykazu to inny rodzaj rozjazdu: nie ślad po tabeli zdjętej, lecz odwzorowanie wniesione przed
migracją, która dopiero tabelę założy, ponieważ definicje komend modułów oddano wraz z miejscem danych,
a scalanie kontraktu poszło jednym przebiegiem, podczas gdy migracje tabel powstają moduł po module
w kroku dobudowy rdzenia. Wykaz sam się sprząta: gdy migracja modułu założy tabelę, sprawdzian wypadnie
niepomyślnie z powodu wiersza, który został, a wpis znika wtedy razem z powodem, dla którego powstał.
## budowa/server/internal/dane/roundtable_tury.go
Numer tury nadaje baza, nie rdzeń: numer powstaje jako największy dotychczasowy w tym oknie plus jeden, wewnątrz transakcji zakładającej wiersz. Gdyby liczył go rdzeń, dwie tury uruchomione w tej samej chwili z dwóch urządzeń tego samego konta dostałyby ten sam numer, a więź jednoznaczności okna i numeru odrzuciłaby drugą. Zero w granicy liczby wypowiedzi znaczy bez granicy, więc jedno przygotowane zapytanie obsługuje zarówno wykaz pełny, jak i wykaz przycięty. Podniesienie numeru redakcji razem z treścią przy zastąpieniu wypowiedzi wynika z tego, że regeneracja jest zastąpieniem, a nie dopisaniem, więc licznik redakcji jest jedynym śladem, że treść się zmieniła. Wskazanie tury nieistniejącej przy dopisywaniu wypowiedzi nie dopisuje niczego i wraca jako ErrBrakWiersza, bo cichy zapis do nieistniejącej tury ukrywałby błąd wołającego. Zapis całej debaty w porządku tur i wypowiedzi istnieje, bo bez niego każda analiza całej debaty czytałaby tury osobno i składała je ręcznie w kolejności.

## budowa/server/internal/podagenci/doradca.go
Pojęcie doradcy trzyma się na trzech własnościach. Jawność rady: Rada niesie
treść wraz z kanałem i modelem doradcy, a Jawnie składa z nich blok do
pokazania; rada nigdy nie wraca samym napisem treści, bo wpleciona w wynik
agenta wyglądałaby na jego własną odpowiedź. Ślad w prowenancji: jedzie
istniejącą drogą rdzenia, a dziennik doradcy zapisuje ten sam opis wywołania
obok pytania i rady, drugiej prowenancji nie ma. Konsultacja, nie delegacja:
metoda Wiazaca oddaje zawsze fałsz, a rama pytania mówi to doradcy wprost,
odpowiedzialność za wynik zostaje przy agencie pytającym.

Kto może być doradcą, rozstrzygają dane, nie kod: czynny wiersz rejestru
kanałów z parametrem doradca i liczbową siłą. Jak daleko wolno sięgnąć,
rozstrzyga sufit siły opisany w doborze doradcy; model z własnej inicjatywy
w górę nie sięga.

Wskazanego doradcy jako osobnego pola tu nie ma, bo Operator dziś nie ma czym
takiego wskazania przyjąć, a takie pole znosiłoby sufit przy każdej
konsultacji i kłamało o powodzie doboru; prośba o silniejszego bez wskazania
Operatora kończy się odmową, nie podmianą.

skrotKonsultacji: doradca i chwila wchodzą do skrótu, bo skrót zdarzenia
konsultacji ma dać się zestawić z jednym wierszem dziennika konsultacji
doradcy; bez nich dwie konsultacje o tym samym pytaniu u dwóch różnych
doradców miałyby skrót identyczny, więc jedyna wartość niosąca treść rady
w zdarzeniu nie wskazywałaby niczego jednoznacznie.

## budowa/server/internal/dane/automations_harmonogram.go

Pole odwołania podpisu niesie wyłącznie referencję klucza HMAC w sejfie poświadczeń, nigdy
jego wartość.

## budowa/server/internal/dane/narzedzia_sesji.go
Narzędzie dołożone komendą po ukośniku trwa do końca sesji: nie wchodzi na
stałe do definicji eksperta i nie znika po jednej turze. Czas życia niesie
klucz obcy z kasowaniem kaskadowym, nie kod — dlatego nie ma tu metody
sprzątającej wygasłe wiersze; koniec życia dołożenia to koniec życia wiersza
sesji.

Plik prowadzi dołożenia jednej sesji, a nie katalog, z którego się je
wybiera. Wykaz pozycji po ukośniku nie jest tabelą: składa go rdzeń na
bieżąco z komend kontraktu, katalogu akcji i katalogu rozszerzeń. Odpisanie
go do trzeciej tabeli byłoby drugą prawdą o tym, co platforma umie,
i rozjechałoby się z pierwszą przy pierwszej instalacji rozszerzenia.

Wiersz niesie odpis, nie odwołanie: trzy pola nazewnicze i grupę, a nie
klucz obcy do pozycji wykazu — wykaz nie jest tabelą i nie ma czego
wskazać. Dzięki temu odinstalowanie rozszerzenia nie zamienia dołożenia
w nazwę bez opisu i do końca sesji widać, co model dostał.

Identyfikator sesji jest tu kluczem wiersza sesji, nie identyfikatorem
kontraktowym: przekład jednego na drugi należy do adaptera rdzenia, tak
samo jak przy każdym innym repozytorium tego pakietu.

Sesja bez dołożeń oddaje wykaz pusty i jest to stan poprawny, nie brak
wiersza — zestaw narzędzi tury jest wtedy samą definicją eksperta.

Powtórzenie dołożenia nie jest błędem i nie mnoży wierszy.

Usunięcie sesji zdejmuje dołożenia kaskadą klucza obcego; zdjęcie
wszystkich dołożeń obsługuje przypadek, w którym sesja zostaje, a zestaw ma
wrócić do podstawy.

Metoda dostępu do repozytorium, a nie pole struktury — z tego samego
powodu, co pozostałe repozytoria pakietu: repozytorium nie trzyma stanu
poza wskaźnikiem na wspólną pamięć zapytań.

Odczyt i zapis dołożenia idą jedną transakcją, bo razem odpowiadają na
jedno pytanie: czy dołożenie już istnieje, a jeśli nie — wpisz je. Bez
wspólnej transakcji dwie komendy po ukośniku wydane w tej samej chwili obie
zastałyby pustą tabelę i obie próbowałyby wpisać wiersz; druga odbiłaby się
o warunek unikalności i dostałaby odmowę za czynność, która była poprawna.

Czas dołożenia ma mówić, kiedy model dostał narzędzie — odświeżony przy
każdym powtórzeniu kłamałby o chwili, od której je ma.

## budowa/server/internal/store/skutek_dobudowy_studia_test.go

Migracja, która przeszła, nie jest dowodem: krok wykonany bez błędu zostawia wersję w dzienniku migracji
niezależnie od tego, czy polecenie czegokolwiek dokonało. Dlatego sprawdzian nie pyta o wersję schematu,
tylko o tabele, kolumny i wiersz katalogu okien, czyli o to, po co ta migracja powstała. Sama definicja
wiersza katalogu okien bez przypięcia do modułu zostawiłaby okno poza zakresem modułu, czyli dokładnie
tam, gdzie było przed migracją.

## budowa/server/internal/dane/automations_kanwa.go

Notatka i położenie węzła kroku zapisują się osobno, więc każdy zapis dotyka wyłącznie
swoich kolumn: ustawienie notatki nie przesuwa węzła, a przesunięcie węzła nie kasuje notatki.

## budowa/server/internal/dane/narzedzia_zakresy.go
Zakres jest nastawą zasięgu, nie bramką wbudowaną: brak wiersza znaczy pełny
dostęp bez granicy, a wiersz powstaje dopiero wtedy, gdy operator coś
przestawił. Tak samo brzmi kontrakt ustawiania zakresu narzędzia.

Zużycie liczy się wierszami, nie licznikiem. Limit obowiązuje w oknie
czasu, więc licznik narastający musiałby być zerowany przez coś, co wie,
kiedy okno się przesunęło. Wiersze ze znacznikiem czasu odpowiadają na
pytanie o liczbę wywołań w ostatnich sekundach jednym zapytaniem i nie
wymagają niczego w tle. Wiersze starsze niż okno sprzątane są przy zapisie
kolejnego — tabela nie rośnie w nieskończoność, a sprzątanie nie potrzebuje
budzika.

Kod sesji pusty przy odczycie zużycia liczy wywołania wszystkich kart.

## budowa/server/internal/podagenci/doradca_wybor.go
Doradca zostaje, konsultacja zostaje, znika wyłącznie ruch modelu w górę; ruch
w górę robi Operator albo nikt. Dwie reguły doboru w ustalonej kolejności:
prośba modelu podlega sufitowi — kanał, o który prosi sam wołający, przechodzi
tylko wtedy, gdy jest dopuszczony do radzenia i gdy da się wykazać, że nie
jest silniejszy od kanału pytającego, a prośba, której nie da się wykazać,
kończy się odmową opisującą brak, nie cichym podstawieniem słabszego; bez
prośby rdzeń bierze najsilniejszego kandydata o sile nie wyższej niż
pytający, a gdy takiego nie ma, kanał samego pytającego — sufit zwężony do
równości jest zawsze wykonalny, więc druga reguła nie potrzebuje odmowy.

Reguły „wskazanie Operatora znosi sufit” tu nie ma: produkt nie ma ani pola
doradcy okna, ani komendy, którą Operator by doradcę wskazał; kanał modelu
okna mówi, którym modelem pracuje okno, a nie kto jest doradcą, więc sufitu
nie znosi.

Porządek modeli pochodzi z danych, nie z kodu: jedyną miarą siły, jaką
produkt ma, jest parametr siły wiersza rejestru kanałów. Nie ma w produkcie
ani katalogu modeli z rangami, ani pola rangi w kontrakcie, ani zestawu
początkowego, który by siłę wypełniał, więc siła niewpisana jest nieznana,
nie zerowa.

Siła nieznana znaczy: nie ma czym zmierzyć, więc w górę się nie idzie
i kandydatem taki kanał nie jest. Zero wpuszczałoby kanał bez siły, a także
literówkę czy wartość ułamkową w parametrze, wszędzie, otwierając sufit
w obie strony; stan produkcyjny, dopóki nikt nie wpisał siły, daje wtedy
kanał pytającego, czyli ten sam model — dokładnie tyle, ile porządek z danych
pozwala orzec.

Kanał samego pytającego w podSufit przechodzi bez mierzenia jako jedyne
odstępstwo: równość z samym sobą zachodzi z definicji, więc nie ma czego
wykazywać nawet wtedy, gdy siła nie jest wpisana; model, który prosi wprost
o swój własny kanał, dostaje w ten sposób to samo, co dostałby bez prośby.

Kanał pytającego w podSufitem jest zapasem, nie wykluczeniem: model o takich
samych parametrach jest dozwolonym rozmówcą, a przy braku innych kandydatów
jedynym możliwym; konsultacja u modelu równego nadal ma sens, bo doradca
dostaje ramę konsultacji i czyste pytanie zamiast całej historii tury. Próg
nieznany nie wpuszcza nikogo obcego, bo wtedy o żadnym kandydacie nie da się
orzec, że nie jest silniejszy.

Sprawdzenie parametru doradca w dopuszczonyZWykazu pilnuje, żeby prośba
modelu nie omijała wykazu kandydatów i nie sięgała po dowolny czynny kanał
rejestru, także taki, którego Operator do radzenia nie dopuścił. Wiersz
wyłączony jest tu tym samym co nieistniejący: kanał, który Operator zgasił,
nie odpowie. Kanał samego pytającego przechodzi bez tego parametru, bo pyta
wtedy sam siebie, a na to nie potrzeba dopuszczenia, którego druga reguła też
nie wymaga.

## budowa/server/internal/store/spojnosc_test.go

Sprawdzian musi dowieść nie tylko tego, że kontrola przechodzi na bazie zdrowej, co pokazuje osobny
sprawdzian przejazdu migracji, lecz przede wszystkim tego, że na bazie chorej nie przechodzi. Kontrola,
która nigdy nie odmawia, jest gorsza niż jej brak, bo daje spokój, którego nie ma czym pokryć. Sierota
w sprawdzianie naruszenia klucza obcego powstaje z połączenia obocznego z wyłączoną pragmą kluczy
obcych — inaczej się nie da, ponieważ pula rdzenia trzyma pragmę włączoną w DSN, więc każdy zapis tą
drogą zostałby odrzucony przy wstawianiu, czyli sprawdzian mierzyłby pragmę, a nie samą kontrolę.
## budowa/server/internal/dane/rozmowa.go
Utrwalany łańcuch przechodzi przez środowisko, moduł, kartę sesji, sesję i okno komunikacji aż do wiadomości. Warstwa danych nie zna pakietu sesji ani rdzenia: jedynym stykiem jest opis okna wraz z funkcją, która go podaje, a rdzeń wypełnia ją swoim rejestrem okien przy montażu. Pole niosące identyfikator okna źródłowego jest identyfikatorem rdzenia w postaci napisu, nie kluczem wiersza, bo warstwa wyższa kluczy wierszy okien nie zna, a przekład na klucz obcy wykonuje utrwalacz przez ten sam łańcuch, co dla okna samej wiadomości. Warstwa wyższa wypełnia metadane przy nadaniu wiadomości, utrwalacz przenosi wartości do kolumn tabeli, a przy odczycie odtwarza je z tych samych kolumn, tak że zapis i odczyt dają tę samą treść. Załączniki jadą w osobnej kolumnie tego samego wiersza: bez nich model wracający do rozmowy nie wiedziałby, że w niej były pliki.

## budowa/server/internal/podagenci/dziennik_doradcy.go
Trwałość konsultacji stoi w tym pakiecie, nie w pakiecie dane, tak jak
dziennik transkrypcji mowy: pojęcie rady nie zmienia się, gdy zmienia się
tabela, a tabela nie zmienia się, gdy przestawia się zasady doboru doradcy.

## budowa/server/internal/dane/okna_operacyjne.go
Okno operacyjne to pozycja katalogu funkcji modułu — nie jest oknem
komunikacji. Okno komunikacji jest bytem wykonania jednej sesji i mieszka
w osobnej tabeli; okno operacyjne jest wpisem rejestru mówiącym, jakie okna
robocze niesie moduł. Wiersze wnosi zaczyn schematu — repozytorium ich nie
zakłada.

Rejestr jest wykazem informacyjnym, nie bramą: moduł nieznany daje wykaz
pusty, nie błąd.

## budowa/server/internal/dane/automations_nadzor.go

Odczyt bazy nie może wynieść wartości poświadczenia, bo kolumny na wartość nie ma — nie
dlatego, że ktoś pamiętał o filtrze przy odczycie. Wartość leży wyłącznie w sejfie plikowym
katalogu danych, tym samym, którym jadą sekrety kont i punktów dostępu.

Zero w miejscu automatyki w wykazie reguł alarmowania znaczy reguły wszystkich automatyk.

Zakres dat w audycie pusty znaczy brak zawężenia: znacznik pusty jest leksykograficznie
mniejszy od każdego znacznika ISO, a górna granica pusta zdejmuje warunek jawnym porównaniem.

## budowa/server/internal/dane/orchestration.go
Podmiana całego układu zależności ma sens tam, gdzie generator oddaje pełną
definicję automatyki naraz. Tu operacja obejmuje jeden łuk: podmiana
kompletu przy dołożeniu pojedynczego łuku przepisywałaby pozostałe wiersze,
gubiąc ich datę utworzenia, a dwa okna pracujące równocześnie kasowałyby
sobie zmiany nawzajem. Obie drogi sięgają tych samych wierszy tej samej
tabeli.

Repozytorium nie ocenia układu: łuk do kroku nieistniejącego i łuk
domykający cykl zapisują się tak samo jak każdy inny — ocena należy do
walidacji w rdzeniu. Więzy pilnowane przez sam schemat, czyli pętla własna
i rodzaj spoza wartości kontraktu, wracają stąd jako błąd zapisu.
## budowa/server/internal/dane/rozmowa_lancuch.go
Reguła całego pliku: brakujące ogniwo łańcucha zakładamy albo zastępujemy najbliższym sensownym, zamiast odmawiać zapisu z powodu braku konfiguracji; odmowa zostaje wyłącznie tam, gdzie nie ma z czego zbudować wiersza. Moduł zastępczy jest jedynym modułem widocznym we wszystkich środowiskach udostępniających moduły, więc zastępstwo nie wprowadza okna do środowiska, w którym nie może się pojawić. Sesja istniejąca wraca bez zmian, bo czynność zapewnienia sesji jest idempotentna: karta sesji i środowisko powstają tą samą drogą co przy utrwalaniu rozmowy, osobnego łańcucha tu nie ma.

## budowa/server/internal/store/zrodlo_migracji.go

Wykaz plików migracji w tym katalogu niesie mniej plików, niż wskazuje najwyższy numer. Wolne numery
nie oznaczają kroku usuniętego ani zgubionego: powstają przy pracy równoległej, gdy numer zarezerwowany
z góry dla kroku, który ostatecznie nie wszedł, zostaje pusty. Numeru zwolnionego nie wolno użyć
powtórnie, ponieważ bazy założone wcześniej mają już wyższą wersję schematu i krok wstawiony w lukę
nigdy by się na nich nie wykonał.

Od numeracji naprawdę wymagane są trzy własności. Jednoznaczność: dwa pliki o tym samym numerze są
awarią startu, bo rejestr migracji ma na kolumnie wersja warunek unikalności, więc drugi krok wykonałby
swój schemat, ale nie zostałby odnotowany. Porządek rosnący: kroki stosuje się po numerze rosnąco, więc
kolejność plików w katalogu nie ma znaczenia. Niezmienność treści: suma kontrolna kroku już zastosowanego
musi się zgadzać przy każdym kolejnym uruchomieniu migracji. Ciągłość numeracji nie jest wymagana.

## budowa/server/internal/podagenci/narzedzia_modelu.go
Wykaz narzędzi modelu powstaje wyłącznie z sekcji pozycji narzędzi kontraktu,
z której generator wytwarza odwzorowanie nazw i komend narzędzi. Serwer
narzędzi czyta ten jeden wykaz: dopisanie komendy do sekcji narzędzi kontraktu
powiększa ten serwer bez zmiany choćby jednej linii kodu. Wpis danaco
w konfiguracji MCP okna powstaje osobno dla każdego okna, z jego
identyfikatorem.

Wynikają z tego dwie rzeczy. Nie ma drugiego miejsca rejestracji narzędzia:
jedyną drogą jest pozycja w sekcji narzędzi kontraktu, a plik kontraktu jest
plikiem zakazanym dla tego pakietu, więc dopisanie pozycji idzie zgłoszeniem
wpięcia, którego dokładną treść niesie PozycjeWpiecia. Budowanie tu własnego
wykazu narzędzi byłoby drugą prawdą o wykazie, dlatego ten plik wyłącznie
czyta odwzorowanie komend narzędzi i odpowiada, czy komendy rodziny
subagent.* już w nim są.

Nazwa spoza wykazu nie znika po cichu: granica uprawnień modelu oddaje
modelowi czytelną odmowę, że komenda stoi poza wykazem narzędzi, i przesuwa
się w chwili, w której przesunie ją kontrakt.

Zdania zastosowania w PozycjeWpiecia mówią modelowi, kiedy sięgnąć po
narzędzie, wzorem pozycji istniejących: powołanie oznacza, że podagent jest
zadaniem w tle pod oknem wykonawcy, powoływanym w trakcie tury.

## budowa/server/internal/dane/automations_okna_wykonania.go

Zero w miejscu automatyki albo harmonogramu w wykazie wyzwoleń znaczy brak zawężenia po tym
polu.

## budowa/server/internal/podagenci/rozruch.go
Powołanie kilkunastu podagentów jednym wywołaniem puszcza tyleż działań naraz;
każde zakłada kolejkę, dokłada pozycję, wiąże ją z wierszem i przestawia stan.
Baza stoi w trybie WAL z limitem czekania na zajętość: zwykły zapis równoległy
nie zawodzi, ale transakcja, która najpierw czyta, a potem pisze, zawodzi, bo
podniesienie blokady odczytu do zapisu nie jest objęte tym limitem i wraca
natychmiast błędem zajętości bazy. Bramka szereguje te transakcje, więc nie
rywalizują o blokadę i nie zawodzą.

Bramka obejmuje wyłącznie założenie pracy, kilka zapisów trwających
milisekundy. Samej pracy podagenta, tury modelu trwającej sekundy albo
minuty, nie obejmuje i obejmować nie może: podagenci mają pracować
równolegle, a bramka rozciągnięta na turę zamieniłaby sieć kilkunastu w
gęsiego idącą jedynkę. Wołający wchodzi w bramkę przed pierwszym zapisem
i wychodzi z niej przed wywołaniem silnika.

Drugie miejsce rozruchu nie dokłada przepustowości, dokłada rywalizację,
czyli dokładnie to, co ta bramka usuwa; wartość liczby miejsc jest stałą, nie
nastawą, bo nastawa bez pytania, które by ją rozstrzygało, byłaby pokrętłem
bez skali.

Liczba miejsc poniżej jedynki w NowaBramka znaczy jedno miejsce: bramka o
zerze miejsc nie wpuściłaby nikogo nigdy, czyli byłaby zatrzymaniem platformy
pod nazwą przepustowości.

Zwolnienie miejsca w Wpusc oddaje się zawsze, także po błędzie założenia
pracy; miejsce niezwrócone zabrałoby sieci przepustowość na stałe. Podagent
odwołany w kolejce do bramki nie ma po co dostać miejsca, więc wynik jest
wtedy fałszem, a zwolnienie mimo to wolno wywołać, bo nic nie robi — dzięki
temu wołający nie musi rozgałęziać odroczonego wywołania. Bramka pusta
wpuszcza natychmiast.

## budowa/server/internal/store/zrodlo_migracji_test.go

Przejazd migracji musi zatrzymać się na kroku o nazwie spoza wzorca, a nie pominąć go po cichu, ponieważ
pominięty krok to schemat niepełny bez jednego komunikatu o błędzie. Numeru zwolnionego nie wolno użyć
powtórnie, więc sprawdzian, który pilnowałby ciągłości numeracji, wymuszałby błąd zamiast go łapać.

## budowa/server/internal/dane/okna.go
Kod eksperta, nie identyfikator wiersza — ekspert bywa kasowany niezależnie
od okien, w których pracował.

Bez identyfikatora zewnętrznego po restarcie rdzenia nie da się połączyć
okna wskazanego przez klienta z jego historią.

Zapis rozmowy jest wołany po każdej turze, bo program wiersza poleceń może
nadać identyfikator dopiero w trakcie pierwszej wymiany. Zapis pustego
napisu jest dozwolony i znaczy rozpoczęcie następnej tury od nowa —
na przykład po przeniesieniu kontekstu.
## budowa/server/internal/dane/schowek.go
Rdzeń nie czyta schowka maszyny Operatora i ta warstwa niczego takiego nie udaje: treść przychodzi z okna, które ją skopiowało, a wraca do okna, które ma ją wkleić, więc repozytorium daje historii wyłącznie trwałość. Powtórzenie treści nie mnoży wierszy: odcisk treści ma warunek jednoznaczności, więc drugi zapis tej samej treści podnosi wiersz zastany na czoło wykazu jednym zapytaniem, bez odczytu i zapisu, który dwa okna kopiujące naraz umiałyby rozjechać. Wpisu bez zapisanej postaci nie wolno czytać jako wpisu o postaci pustej: puste znaczy, że nie wiadomo, jaka była postać źródła, więc wklejenie z poleceniem zachowania postaci źródła ma wtedy zachować postać miejsca docelowego. Postać fragmentu nie wchodzi do rachunku odcisku, bo odcisk odpowiada na pytanie, czy tę samą treść już odłożono, a ten sam akapit skopiowany dwa razy w różnym wyróżnieniu ma zostać jednym wpisem historii.

## budowa/server/internal/dane/automations_szablony.go

ZapiszSzablonAutomatyki zapisuje szablon i jego parametry jedną transakcją, bo szablon
zapisany z parametrami poprzedniej wersji byłby formularzem pytającym o pola, których
szablon już nie zna.

## budowa/server/internal/podagenci/siec_granica.go
Przycinanie samego wywołania przepuszcza dwa wywołania po piętnaście, bo
granica dotyczy stanu całej sieci okna, nie kształtu pojedynczego żądania.
Miejsce w sieci zajmują wyłącznie podagenci czynni, czyli ci, którzy pracę
mają przed sobą albo w toku; podagent zakończony, błędny i zatrzymany
miejsca nie trzyma, inaczej okno zużyłoby piętnaście miejsc raz na całą swoją
historię.

Żądanie większe niż liczba wolnych miejsc wykonuje się na tylu, ile jest
wolnych. Odmowa zostaje na jeden przypadek: sieć pełna, gdzie przyciąć można
wyłącznie do zera, a powołanie zerowe udawałoby wykonanie; odmowa niesie
wtedy zdanie mówiące, co zrobić — zebrać wyniki albo zatrzymać któregoś.

Orzeczenie stanu końcowego w StanKoncowy jest wspólne dla granicy sieci, kto
zajmuje miejsce, zbierania wyników, czy komplet gotowy, i domknięcia
zatrzymania, czy wiersz wciąż kłamie; trzy osobne odpowiedzi rozjechałyby się
przy pierwszej zmianie słownika stanów.

Wiersze podawane do MiejscaWSieci muszą pochodzić z jednego okna: granica
jest granicą sieci wykonawcy, nie granicą platformy — liczenie kompletu bazy
zamknęłoby powołanie u wszystkich, gdy jedno okno zapełni swoją sieć.

Bez repozytorium albo bez kodu okna granica w PrzydzialOkna nie ma czego
liczyć i oddaje przydział samego wywołania; przydział wyliczony z nieudanego
odczytu byłby granicą zmyśloną.

## budowa/server/internal/dane/automations_wersje.go

Zapis wersji jest idempotentny po parze automatyki i numeru: ten sam zapis definicji
powtórzony nie zakłada drugiej migawki tego samego numeru.

Publikacja jest jedna na automatykę: zdjęcie znacznika publikacji ze wszystkich wersji
i nadanie go jednej idzie w tej samej transakcji, żeby historia nie pokazała przez chwilę
dwóch wersji opublikowanych naraz.
## budowa/server/internal/dane/sejf_poswiadczen.go
Baza zna wyłącznie odwołanie, czyli nazwę wpisu w sejfie; sam sekret leży w osobnym pliku katalogu danych. Dzięki temu odczyt katalogu kont nigdy nie może wynieść sekretu, bo kolumny na niego nie ma, a mimo to poświadczenie jest trwałe, więc znacznik posiadania poświadczenia mówi prawdę: skoro odwołanie zapisano, sekret istnieje. Sejf jest celowo prosty: jeden plik JSON kluczowany bytem, pod zamkiem. To nie jest magazyn klasy zarządzania kluczami, tylko trwały schowek na sekret, którego rdzeń nie wpuszcza do bazy ani do odpowiedzi; wymianę na zewnętrzny magazyn domyka ten sam interfejs zapisu i usuwania. Byt w odczycie jest surowym kluczem wpisu, bez przedrostka odwołania: rozbiera go wołający, a sejf kluczuje bytem dokładnie tak, jak zapisał przy zapisie poświadczenia.

## budowa/server/internal/tokenizator/tokenizator.go

Przybliżenie liczby żetonów po długości tekstu myli się o kilkadziesiąt procent na języku polskim,
ponieważ znaki diakrytyczne rozpadają się na osobne żetony, a jeszcze bardziej na kodzie źródłowym i na
zapisie strukturalnym. Licznik, który myli się o kilkadziesiąt procent, jest gorszy niż jego brak: pasek
zajętości pokazuje stan w normie, a tura kończy się przepełnieniem okna kontekstu. Słownik BPE jest
wkompilowany w binarium rdzenia razem z pakietem — nie ma tu ani jednego procesu potomnego, ani jednego
pobrania z sieci w trakcie żądania, ponieważ tokenizator ma działać na maszynie odciętej od świata tak
samo, jak na maszynie deweloperskiej.

Pakiet nie udaje, że zna podział na żetony każdego istniejącego modelu. Zna rodziny słowników, które
naprawdę ma; model spoza nich dostaje słownik najbliższy wraz z jawnym powiedzeniem, którym słownikiem
policzono. Nazwa słownika wchodzi do odpowiedzi kontraktu, więc czytelnik wie, czym zmierzono, zamiast
dostać liczbę bez świadka. Wybór słownika domyślnego dla modeli spoza wykazu jest świadomy, a nie
zaniedbaniem: jest najbliższym dostępnym podziałem dla tej klasy modeli, a odpowiedź mówi, że policzono
właśnie nim, więc czytelnik wie, na ile liczba jest wiążąca.

## budowa/server/internal/dane/orkiestracja_biegi.go
Obsada jest jednym bytem o dwóch nośnikach: wisi albo na automatyce, albo na
biegu orkiestracji, pilnowane więzem sprawdzającym schematu. Zapytania
pętli automatyki niosą identyfikator automatyki i widzą wyłącznie obsadę
automatyk; zapytania tego pliku niosą identyfikator biegu i widzą wyłącznie
obsadę biegów — stąd dwa pliki nad jedną tabelą.

Bieg nie jest przebiegiem automatyki: przebieg automatyki wisi na zapisanej
definicji, a bieg orkiestracji zakłada się w oknie i bywa jednorazowy —
zestawienie dwóch modeli nie wymaga zapisanej automatyki.

Kolumna wskazująca kolejkę jest wskazaniem, nie drugim silnikiem — pracę
biegu wykonuje silnik kolejek.

Okno bez biegu wraca jako błąd braku wiersza i jest to stan zwykły, nie
usterka: podagent bywa powołany poza biegiem orkiestracji, w zwykłej
rozmowie — kolumna biegu podagenta dopuszcza pustkę.

Podmiana obsady, a nie dopisywanie: scalanie zostawiałoby stanowiska
usunięte z nadesłanego wykazu. Obsada pusta jest poprawna — bieg bez
obsady rusza na modelu wskazanym w oknie.

## budowa/server/internal/dane/badania.go

Sygnatury dobudowy modułu stoją w badania_dobudowa.go; osadzenie ich przez wbudowany
interfejs trzyma cały obszar w jednym kontrakcie, nie rozbija go na dwa niezależne porty.

## budowa/server/internal/podagenci/zywotnosc.go
Rejestr procesów sesji obejmuje proces tury uchwytem systemowy, grupa
procesów na Uniksie, obiekt zadania na Windows, a dogląd okresowy utrzymuje
odpowiedź o życiu zgodną z prawdą systemu, nie z polem w pamięci. To jest
żywotność realna i ten plik z niej wyłącznie czyta.

Droga dziś nie ma dwóch rzeczy. Tura pozycji kolejki, a praca podagenta jest
pozycją kolejki, jedzie z pustymi zasięgami, więc proces wykonujący pozycję
nigdy nie trafia do rejestru procesów. Rejestr kluczuje procesy
identyfikatorem okna i trzyma jeden wpis na okno, a przejęcie ubija wpis
poprzedni; podagentów bywa piętnastu pod jednym oknem i pracują równolegle
z turą własnego okna wykonawcy, więc zarejestrowanie ich procesów pod oknem
wykonawcy ubijałoby nawzajem turę orkiestratora i tury podagentów. Brakuje
klucza drobniejszego niż okno, i tego pakiet nie obchodzi bokiem, bo drugi
rejestr procesów byłby drugą prawdą o procesach.

Żywy jest wobec tego mierzalnie proces orkiestratora, okna wykonawcy, które
podagentów powołało: jego turę startuje wysłanie wiadomości, zasięg okna
jest wtedy wypełniony i proces zostaje wpisany do rejestru. Ocena mówi więc
prawdę o oknie prowadzącym podagentów, a o procesie samej pozycji mówi stan
bez wpisu, i to zdanie jest prawdziwe, nie zastępcze.

Żywotność po awarii i po restarcie stanowi drugą połowę tego pliku. Dogląd
rejestru mówi o procesie rdzenia, który stoi; gdy rdzeń padnie, nie mówi nic
i nie ma komu mówić, dlatego pytanie, co się dzieje z podagentem po awarii,
ma odpowiedź w bazie, nie w rejestrze procesów.

Znacznik uruchomienia składa się z dwóch rzeczy, bo żadna sama nie
wystarcza: numer procesu jest w systemie powtarzalny, po restarcie maszyny
ten sam numer wraca, a czas sam nie odróżnia dwóch rdzeni wstałych w tej
samej milisekundzie. Razem są jednoznaczne w praktyce, a jednoznaczności
absolutnej ten znacznik nie potrzebuje: rozstrzyga wyłącznie pytanie, czy to
nadal ten sam rdzeń. Znacznik zakłada się raz na proces i podaje dalej
wartością, losowania po drodze nie ma, więc nikt nie osieroci sam siebie.

PosprzatajPoRestarcie woła się raz, przy starcie, przed pierwszym powołaniem:
wywołanie późniejsze zamknęłoby pracę powołaną przez ten sam rdzeń, gdyby jej
oznaczenie prowadzenia jeszcze nie doszło. Trwałość pusta znosi się sama:
rdzeń bez repozytorium podagentów startuje, a nie odmawia startu, po prostu
nie ma czego sprzątać.

## budowa/server/internal/dane/orkiestracja_podagenci_zywotnosc.go
Osobny plik, a nie dopisek do pliku trwałości powołania: tamten plik
odpowiada za trwałość powołania, stanu i wyniku — byt, który ma restart
przeżyć. Tutaj stoi rzecz przeciwna: praca, która restartu przeżyć nie
może, bo żyje w procesie modelu. Dwie odpowiedzialności, dwa pliki.

Skąd się biorą sieroty: podagent w stanie oczekującym albo działającym
opisuje pracę, którą ktoś wykonuje. Wykonawcą jest proces modelu wystawiony
przez rdzeń, więc rdzeń ubity zabiera go ze sobą, a wiersz zostaje z zapisem
nieprawdziwym. Rozstrzyga o tym znacznik uruchomienia: wiersz niezakończony
prowadzony przez uruchomienie inne niż bieżące jest sierotą, bo jego
wykonawcy nie ma.

Zamknięcie, nie wskrzeszenie: sierota idzie w stan zakończony błędem
z powodem osierocenia, a nie z powrotem w stan oczekujący. Powtórne
puszczenie pracy byłoby decyzją, której nikt nie podjął — operator ma
zobaczyć, co się urwało, i powołać na nowo sam. Wynik już zebrany zostaje:
podagent, który zdążył coś oddać przed awarią, oddaje to nadal.

Wykaz podagentów pusty przy oznaczaniu prowadzenia nie jest błędem.

Zamknięcie sierot woła się raz, przy starcie rdzenia, zanim jakikolwiek
podagent zostanie powołany, i oddaje wykaz zamkniętych do meldunku, bo
cicha zmiana stanu na wykazie operatora byłaby zmianą niewidoczną.

Tym samym powodem, co przy powołaniu podagentów: połowa powołania
oznaczona, a połowa nie, dałaby przy następnym starcie sieroty z podagentów
właśnie pracujących.

Odczyt i zapis zamknięcia sierot idą w jednej transakcji, bo między nimi
nie ma prawa wejść powołanie nowego podagenta — zamknęłoby się dopiero co
powołaną pracę. Sierot brak jest odpowiedzią poprawną i najczęstszą: rdzeń
zamknięty porządnie zostawia wszystkich w stanie końcowym.
## budowa/server/internal/dane/sesje.go
Sesje w koszu nie wchodzą do wykazu sesji żywych: widzi je wyłącznie repozytorium kosza, a odczyt po identyfikatorze ich nie kryje, bo po nim odbywa się przywrócenie i czyszczenie. Usunięcie sesji przez tę warstwę kasuje sesję wraz z oknami i wiadomościami kaskadą schematu bazy, ale jest to prymityw warstwy danych: torem produktu jest kosz, gdzie polecenie usunięcia stawia znacznik, a fizyczne skasowanie wykonuje czyszczenie po terminie, bo tylko ono sprząta też bloki wiadomości.

## budowa/server/internal/protocol/fragment.go
Numeru fragmentu ani znacznika końca w typie Chunk nie ma: kontrakt umieszcza
je w kopercie, polami sekwencji i zakończenia, wstawiane przy złożeniu koperty
fragmentu i odczytywane przy odczycie numeru i ostatniego fragmentu.

Odbiorca zastępuje tekstem wersji ostatecznej tekst złożony z fragmentów
tekstowych, zamiast dopisywać go na końcu; rodzaj wersji ostatecznej odróżnia
całość od ciągu dalszego, a fragment pakuje się kopertą ze znacznikiem końca,
jako zdarzenie domykające strumień.

ChunkBledu niesie błąd techniczny wywołania, a przyczyna jedzie i jako treść
nietekstowa, i jako tekst dla użytkownika; fragment kończy strumień, więc
wywołujący pakuje go kopertą z ostatnim równym prawda, choć sesja i konto
pozostają przy tym czynne.

## budowa/server/internal/transport/bramka.go

Bez straży wejścia klient bez tokenu wykonuje dowolną komendę, w tym polecenie zdalne, które przez tor
zdalny potrafi sięgnąć po połączenie SSH na maszynę Operatora. Zgoda na hosta zdalnego jest wydana
hostowi, nie wołającemu, więc wystawiony rdzeń oddawałby cudzemu połączeniu zarówno zdalny serwer, jak
i komputer w biurze. To nie jest bramkowanie uprawnień, lecz wskazanie miejsca, w którym logowanie do
aplikacji ma skutek dla rdzenia: straż nie zna pojęcia uprawnienia, roli, zakresu ani modułu, pyta tylko,
czy dane gniazdo przeszło przez bramkę; nie pyta o to ani razu na pętli zwrotnej, bo tam jest wyłączona
w całości; po przejściu bramki milczy do końca życia połączenia. Odmowa opisuje brak, nie zakaz — to
połączenie nie przeszło przez bramkę, a droga naprawy w postaci zalogowania się stoi w tym samym zdaniu.

Wykaz komend wejścia jest wyczerpujący i wynika z jednego pytania: czego nie da się pominąć, żeby móc
się zalogować. Powitanie połączenia jest tędy, którędy token wchodzi do rdzenia, ponieważ transportem
jest jedno gniazdo bez nagłówka na każdym żądaniu. Logowanie jest tędy, którędy token powstaje.
Rejestracja jest potrzebna, bo bez niej rdzeń bez założonej bramki byłby zamknięty na klucz, którego
nikt jeszcze nie wykuł. Weryfikacja jest krokiem wydającym token w rejestracji dwukrokowej — bez niego
rejestracja zakłada konto niepotwierdzone, którego już nic nie potwierdzi, ponieważ rejestracja drugi
raz oddaje konflikt, a logowanie odmawia zdaniem o oczekiwaniu na potwierdzenie adresu. Odzyskanie
i zresetowanie konta to dwa kroki naciskane przez tego, kto hasła nie pamięta, czyli z definicji przez
bramkę nie przejdzie, a odbite kroki zamieniają zapomniane hasło w koniec instalacji. Odświeżenie tokenu
przedłuża sesję zapisaną na maszynie: token przedstawiony w powitaniu wiąże gniazdo i wtedy przedłużenie
przechodzi samo, ale token wygasły gniazda nie wiąże, więc Operator ma wtedy zobaczyć odmowę rdzenia
o wygaśnięciu sesji, a nie odmowę straży, która o sesji nic nie mówi. Wykaz nie jest furtką: żadna z tych
komend nie wykonuje pracy Operatora ani nie sięga po pliki, sieć czy powłokę — wszystkie dotykają
wyłącznie bramki, a zgadywanie po nich jest ograniczone dławikiem prób wejścia i tym, że droga
potwierdzenia jest losowa, jednorazowa i wygasa po godzinie. Nazwy komend biorą się ze stałych kontraktu,
nie z literałów, żeby zmiana nazwy komendy w kontrakcie wywróciła kompilację, a nie po cichu zamknęła
wejście.

Wymóg logowania bierze się z dwóch rzeczy: adresu nasłuchu, który poza pętlą zwrotną obowiązuje sam
z siebie jako fakt, a nie nastawa, więc wystawienia nie da się zrobić przez zapomnienie; oraz jawnego
wskazania Operatora, będącego dźwignią włączającą wymóg także na pętli zwrotnej albo znoszącą go przy
nasłuchu szerszym — zniesienie jest dozwolone, ale nigdy ciche, bo dziennik mówi wtedy wprost, co stoi
otworem. Metoda przepusc ma trzy wyjścia na tak i jedno na nie: straż wyłączona przepuszcza wszystko,
komenda wejścia przechodzi zawsze, związane gniazdo przechodzi zawsze, a jedynym przypadkiem odmowy
jest obowiązujący wymóg przy komendzie spoza wejścia i gnieździe nieprzedstawionym. Pytanie dotyczy
gniazda, nie tożsamości wołającego, dlatego przeżyje zmianę modelu bramki: dziś sesja bramki nie ma
właściciela, a gdy dostanie konto i wiele urządzeń z osobnymi tokenami, straż nie będzie wymagać ani
jednej zmiany. Rdzeń nieznający rozszerzenia stanu bramki nie przepuszcza żądania — jedyne miejsce w tym
pakiecie, gdzie brak czegoś zamyka drogę zamiast ją otwierać, ponieważ rdzeń, który nie umie odpowiedzieć,
kto woła, przy obowiązującym wymogu oddawałby komendy komukolwiek.

Kod odmowy bez bramki jest jeden, oznaczający brak uwierzytelnienia z kontraktu, a nie brak uprawnienia
ani błąd pola żądania, ponieważ brakuje właśnie przejścia przez bramkę; klient rozpoznaje ten kod i otwiera
okno logowania zamiast pokazywać błąd komendy.
## budowa/server/internal/dane/sesje_kosz.go
Repozytorium kosza jest osobne od repozytorium sesji, bo tamto obsługuje sesje żywe i jego wykaz sesji z kosza nie widzi; kosz jest odwrotną stroną tej samej tabeli, widzi wyłącznie wiersze ze znacznikiem usunięcia — jedna tabela, dwa pytania. Czyszczenie zabiera też bloki wiadomości: kaskada schematu od sesji sprząta okna i wiadomości, ale bloki wiadomości wiszą na identyfikatorach kontraktowych bez klucza obcego, więc czyszczenie usuwa je wprost, w tej samej transakcji.

## budowa/server/internal/dane/pamiec.go
Treść obszerna trafia do pliku, baza trzyma odwołanie. Konfigurację
pamięci sesji obsługuje osobny plik repozytorium.

## budowa/server/internal/protocol/fragment_test.go
Rozjazd identyfikatora, numeru albo domknięcia zostawia okno w ładowaniu albo
składa treść z dwóch tur naraz.

Fragment niedomykający strumienia nie ma pola znacznika końca wcale, tak
klient odróżnia stan jeszcze nie koniec od stanu koniec równy fałsz.

Sprawdzian ciągłości identyfikatora przechodzi całą turę — żądanie, trzy
fragmenty, domknięcie — i porównuje identyfikatory oraz kolejność numerów.
## budowa/server/internal/dane/skroty_tekstowe.go
Skrót rozwija się we wszystkich polach tekstowych platformy, więc słownik skrótów należy do rdzenia, nie do jednego modułu. Profil pusty znaczy skrót wspólny; zapisuje się pustym napisem, nie wartością pustą bazy, bo warunek jednoznaczności nad kolumną dopuszczającą taką wartość nie pilnowałby niczego.

## budowa/server/internal/dane/biblioteka_kolekcje_pliku.go

UstawKolekcjePliku ma semantykę ustawienia, nie dokładania: stan po zapisie jest wykazem
z żądania, więc zdejmuje plik z kolekcji pominiętych w wykazie. Kolekcja nieznana w żądaniu
jest odmową całości: wykaz z kodem, którego nie ma, wskazuje przynależność nieosiągalną,
a wykonanie reszty zdjęłoby plik z kolekcji zastanych na podstawie żądania zrozumianego
tylko częściowo.

## budowa/server/internal/protocol/koperta.go
Pakiet nie definiuje ani jednej nazwy, kodu, wartości wyliczenia ani kształtu
komunikatu; wszystkie pochodzą z pakietu współdzielonego wytworzonego z pliku
kontraktu, jedynego źródła prawdy. Zmiana kontraktu przerywa kompilację tego
pakietu, zamiast rozjeżdżać się z nim po cichu.

LadunekDo jest funkcją, nie metodą: Koperta jest typem kontraktu, więc
zachowanie dokłada się obok niego, a nie w jego definicji.

## budowa/server/internal/dane/pamiec_konteksty.go
Kontekst jest zestawem wskazań: usunięcie kontekstu nie kasuje ani jednego
wpisu pamięci, bo kontekst nie jest właścicielem treści. Dlatego wpisy idą
zapisem strukturalnym w kolumnie, a nie kluczem obcym z kaskadą.

Kontrakt ustawiania retencji oddaje liczbę wpisów wprost i nie ma prawa
jej zgadywać: pochodzi z policzenia wierszy, nie z oszacowania.

Granica jest znacznikiem czasu w zapisie kolumny daty utworzenia wpisu
pamięci projektu — czyli zapisem w formacie ISO 8601 w strefie uniwersalnej.
Porównanie napisów jest tu poprawne, bo ten zapis rośnie leksykalnie razem
z czasem.

## budowa/server/internal/transport/bramka_test.go

Sprawdziany straży bramki mierzą ją tablicą, bo reguła ma dokładnie trzy wejścia — wymóg, rodzaj
komendy i więź gniazda — i każdy ich układ ma jedno rozstrzygnięcie; pominięcie któregokolwiek pola
tablicy zostawia w straży dziurę wielkości jednej komendy. Dopisanie komendy spoza wykazu komend wejścia
otwiera ją na oścież przed logowaniem, więc zmiana ma się o sprawdzian potknąć. Wykaz liczy siedem
pozycji, nie trzy: rejestracja jest dwukrokowa, bo to weryfikacja wydaje token, nie rejestracja sama,
odzyskanie konta jest dwukrokowe, a przedłużenie sesji dotyczy tokenu, który gniazda nie związał — każda
z nich pada z gniazda jeszcze nieprzedstawionego, i bez każdej z nich któraś droga wejścia jest zamknięta
na głucho.

## budowa/server/internal/protocol/koperta_test.go
Sprawdziany koperty mierzą dokładnie to, co kontrakt obiecuje klientowi:
komunikat niepoprawny strukturalnie jest czymś innym niż komenda nieznana,
a to, co poszło na drut, wraca z drutu bez zmiany.
## budowa/server/internal/dane/slownik.go
Ślady importu i eksportu leżą w pliku sąsiednim, pamięć tłumaczeń w innym pliku sąsiednim: jedno repozytorium rozdzielone na pliki wedle odpowiedzialności, z typem, interfejsem i konstruktorem zadeklarowanymi wyłącznie w jednym z nich. Wystąpienia terminu nie mają tu tabeli: liczą się w locie z treści okna albo panelu przeszukanej względem źródła terminu, żeby nie unieważniać zapisu przy każdej korekcie panelu. Zapis ma jedną drogę: żądanie nadsyła zawsze komplet zmian naraz, a zapis terminów przyjmuje wykaz i zapisuje go w jednej transakcji przez wstawienie z aktualizacją po identyfikatorze zewnętrznym — termin ze wskazanym kodem aktualizuje się, termin bez zastanego wiersza o tym kodzie zakłada się, ten sam zapis obsługuje obie ścieżki, nie dwie osobne metody.

## budowa/server/internal/dane/biblioteka_opis.go

Brak opisu nie jest brakiem wiersza w rozumieniu odmowy: zasób bez ani jednego wypełnionego
pola oddaje opis pusty, nie błąd braku wiersza. Odmowa należy się wyłącznie wskazaniu
zasobu, którego nie ma, co rozstrzyga odczyt pliku, zanim opis w ogóle zostanie odczytany.

Zawężenie definicji pól do rodzaju treści jest miękkie: pole bez wskazania rodzaju treści
stosuje się do wszystkich, więc trafia do wyniku każdego zawężenia — twarde porównanie
zdejmowałoby z formularza pola ogólne przy pierwszym zawężeniu.

Wartość pola niestandardowego leży w zapisie JSON, więc zliczenie zasobów z danym polem
idzie funkcją odczytu JSON silnika bazy, nie po tekście, żeby kod pola będący fragmentem
innego kodu nie dawał fałszywego trafienia.

## budowa/server/internal/protocol/odpowiedz.go
Typ, identyfikator i sesja koperty odpowiedzi pochodzą z żądania, dzięki
czemu klient wiąże odpowiedź z wywołaniem, które ją wywołało. Ścieżka nie
może zawieść, bo wynik jest już zserializowany przy budowie odpowiedzi.

## budowa/server/internal/transport/kontrakt.go

Transport nie zna rdzenia: zna wyłącznie interfejs Rdzen, który rdzeń realizuje, oraz kontrakt
komunikatów. W tym pakiecie nie ma ani jednego importu pakietu wykonawczego — dzięki temu zamiana
rdzenia nie dotyka transportu, a transport da się uruchomić w sprawdzianie z atrapą rdzenia po stronie
sprawdzianu. Serwuje też pliki klienta i utrzymuje wiele połączeń równocześnie.

Do czasu przypisania konta metodą PrzypiszKonto obowiązuje konto z parametru nawiązania połączenia,
ponieważ uwierzytelnianie jest jedyną kontrolą dostępu i w fazie budowy nie działa. Koperta o pustym
polu Type oznacza, że odpowiedź pójdzie osobno — transport nic wtedy nie odsyła, a rdzeń sam korzysta
z ujścia albo rozgłośnika. Kontekst przekazywany rdzeniowi jest kontekstem serwera, nie połączenia:
rozłączenie klienta nie przerywa rozpoczętej pracy rdzenia. Rdzeń opisuje drogę rozgłaszania portem
własnym i sięga wprost po metodę Rozglos serwera — interfejs po stronie transportu byłby drugą
deklaracją tej samej rzeczy.

## budowa/server/internal/dane/biblioteka_slownik.go

Wykaz słownika łączy dwa źródła — wpisy słownika i etykiety nadane przy zasobach — sumą
zdejmującą powtórzenia, więc etykieta obecna w obu miejscach wychodzi w wykazie tylko raz.

EtykietaSlownika: etykieta nosząca zasoby, ale bez własnego wpisu słownikowego, wraca
z licznikiem użycia i pustym czasem założenia, bo istnieje mimo braku wiersza w słowniku.

PrzemianujEtykiete zapisuje przez zastąpienie zamiast zwykłej aktualizacji, bo zasób noszący
obie nazwy naraz złamałby klucz główny pary pliku i etykiety: wiersz stary ustępuje wtedy
nowemu zamiast wywracać całą zmianę.
## budowa/server/internal/dane/slownik_pamiec.go
Zapis pamięci dokłada wiersz bez czytania bieżącej treści panelu w edycji: to wywołujący rozstrzyga, kiedy para segmentów jest zatwierdzona i warta zapamiętania. Dopasowanie podpowiedzi jest przybliżone tylko na tyle, na ile pozwala baza: kontrakt chce dopasowania po podobieństwie, nie po równości, a silnik bazy bez rozszerzenia nie ma wbudowanej miary podobieństwa napisów, więc wyszukiwanie podpowiedzi wykonuje dopasowanie podciągu segmentu źródłowego, nie dopasowanie znaczeniowe ani odległość edycyjną. Czas jest liczbą milisekund epoki, tym samym wzorem co w pozostałych tabelach czasowych repozytorium.

## budowa/server/internal/dane/pamiec_wylaczenia.go
Wyłączenie nie jest usunięciem: ten plik nie ma ani jednego zapytania
dotykającego wpisów pamięci projektu. Wpis wyłączony zostaje na miejscu
z nietkniętą treścią i wraca do kontekstu po usunięciu wiersza wyłączenia.
Tym różni się wyłączenie pamięci od jej usunięcia i tym różni się ten plik
od pliku pamięci obszaru roboczego.

Wyłączenie nie jest też odpięciem. Odpięcie zwęża zasięg samego wpisu,
więc zmienia wiersz pamięci; wyłączenie zostawia wpis nietknięty
i wstrzymuje go wyłącznie w zasięgu, w którym operator go nie chce.

Dwa byty wyłączane: pojedynczy wpis albo cały poziom pamięci. Schemat
pilnuje, żeby wiersz wskazywał dokładnie jeden z nich więzem
sprawdzającym; tutaj pilnuje tego osobna funkcja sprawdzająca byt
wyłączenia, żeby odmowa doszła do operatora z nazwą braku, a nie jako
naruszenie warunku bazy.

Wyłączenie tego samego bytu w tym samym zasięgu, zapisane powtórnie, nie
jest zmianą i nie ma czego rozgłaszać — wyłączenie dwa razy nie jest stanem,
który da się znieść jednym ruchem.

Podzapytanie wskazanego wpisu pamięci oddało wtedy pustą wartość, a warunek
schematu przestawił wiersz na wyłączenie poziomu. Milczenie byłoby tu ciszą
udającą zapis.

## budowa/server/internal/protocol/rozpoznanie.go
Nazwy komend i zdarzeń wstrzykuje punkt wejścia, biorąc je wyłącznie ze
stałych wytworzonych do pakietu shared, wywołaniem budującym rejestr z pełnym
wykazem komend kontraktu. Dzięki temu warstwa protokołu pozostaje wolna od
powielonych literałów, a zbiór nazw znanych rdzeniowi zmienia się wyłącznie
razem z kontraktem.

## budowa/server/internal/dane/biblioteka_sugestie.go

RozstrzygnijSugestie: sugestia już rozstrzygnięta nie liczy się po raz drugi — warunek
zapytania pilnuje tego zamiast wołającego.
## budowa/server/internal/dane/slownik_wymiana.go
Tabele są dwie, bo import i eksport to różne kierunki z różną kolumną wyniku; wspólna tabela byłaby dwiema prawdami o jednym bycie. Ślad eksportu nie niesie dowodu powstania pliku: rdzeń nie ma magazynu blobów, a kontrakt oddaje wyłącznie liczbę wyeksportowanych pozycji, bez identyfikatora pliku ani rozmiaru, więc pole ścieżki niesie ścieżkę żądaną przy wywołaniu, nie ścieżkę wyniku eksportu.

## budowa/server/internal/protocol/tozsamosc_zadania.go
Kontrakt stanowi, że odpowiedź i fragmenty strumienia powtarzają identyfikator
zadania: koperta fragmentu ma nieść ten sam identyfikator, co koperta żądania,
które turę otworzyło. Nadawca strumienia jednak żądania nie widzi, obsługiwacz
komendy dostaje wyłącznie rozpakowany ładunek, a identyfikator zostaje
w żądaniu; bez tego wpisu w kontekście fragmenty tury i odpowiedź niosłyby
różne identyfikatory.

Kontekst niesie tożsamość żądania zamiast nowego parametru, bo droga
alternatywna to poszerzenie podpisu każdej czynności domenowej o identyfikator
żądania — sto z górą podpisów zmienionych po to, by kilka z nich go użyło.
Kontekst niesie już zasięg wykonania i odwołanie tury, więc tożsamość żądania
jest tu bytem tej samej klasy i wchodzi jednym wpisem, wspólnym dla wszystkich
komend.

Klucz wpisu jest typem prywatnym, żeby wykluczyć kolizję z kluczem innego
pakietu: nikt spoza tego pakietu nie ma jak zapisać ani nadpisać tego wpisu
inaczej niż funkcją tego pliku.

Brak wpisu nie jest błędem: tura powołana poza drogą komendy nie ma żądania
i dostaje napis pusty, a wywołujący podstawia wtedy własną tożsamość
zastępczą i mówi o tym wprost, zamiast udawać żądanie, którego nie było.

TozsamoscStrumienia istnieje po to, żeby to rozstrzygnięcie stało w jednym
miejscu: gdyby każdy nadawca strumienia wybierał sam, obietnica jednego
identyfikatora przez cały strumień byłaby powtarzana w kilku plikach,
a obietnica powtórzona to obietnica, którą któryś z nich kiedyś złamie.

## budowa/server/internal/dane/bloki_wiadomosci.go

Tekstu tutaj nie ma: fragment rodzaju text domyka dziennik rozmowy w kolumnie treści
wiadomości; zapisanie go drugi raz tutaj byłoby drugą prawdą o tej samej wypowiedzi, więc
warunek CHECK schematu odbija taki zapis, a rejestrator bloków nawet go nie próbuje.

OknoKod i WiadomoscKod są identyfikatorami kontraktowymi, nie kluczami obcymi, bo rejestrator
strumienia innych repozytoriów nie zna.

Kolejność bloku w obrębie wiadomości wyliczana jest w tym samym poleceniu zapisu, wzorem
numeracji historii okna, bez osobnego odczytu i transakcji.

## budowa/server/internal/transport/nawiazanie.go

Sprawdzanie pochodzenia w metodzie nawiaz nie jest bramką kontrolną: wykaz pochodzeń nie pyta, kim jest
wołający i czego mu wolno, tylko odcina wyłącznie stronę trzecią, która namówiła przeglądarkę operatora
maszyny, żeby otworzyła gniazdo do jego rdzenia; operator nie widzi tego nigdy, widzi to wyłącznie cudza
strona. Klientem bywa webview powłoki, który przedstawia się rozmaitym Origin — pochodzenia własne
obejmują schemat powłoki oraz localhost i adres pętli zwrotnej z dowolnym portem, więc powłoka wchodzi
bez wskazywania czegokolwiek. Wystawienie pod inną domenę dopisuje ją do wykazu pochodzeń dozwolonych.

Wzorce pochodzeń własnych obejmują powłokę aplikacji oraz interfejs otwarty w przeglądarce pod adresem
pętli zwrotnej, a port bywa dowolny, ponieważ nasłuch potrafi wziąć port wskazany przez system — stąd
gwiazdka w porcie, a nie w całym wzorcu. Pętla zwrotna ma dwa adresy, nie jeden: przeglądarka na
maszynie z pierwszeństwem IPv6 rozwiązuje localhost na adres skrócony i podaje wtedy pochodzenie tej
samej pętli zwrotnej, na której rdzeń nasłuchuje, a bez tych wzorców nawiązanie kończyłoby się odmową.
Ukośniki odwrotne w tych wzorcach są konieczne, nie ozdobne: dopasowanie idzie przez dopasowanie ścieżek,
gdzie nawias kwadratowy otwiera klasę znaków, a wzorzec z nawiasem gołym jest wzorcem wadliwym, przy
którym biblioteka gniazda przerywa przegląd wykazu błędem — jeden zły wzorzec potrafi więc odciąć
pochodzenia sprawdzane po nim.

Wykaz pochodzeń wskazany przy wystawieniu dopisuje się do wzorców własnych, nie zastępuje ich, ponieważ
wystawienie pod domenę nie jest powodem, żeby produkt przestał wpuszczać własną powłokę, a taki właśnie
byłby skutek zastąpienia wykazu.

Adres urządzenia w linii dziennika istnieje po to, że telemetria połączeń jedzie dziennikiem rdzenia
i trwale trafia do wpisu diagnostycznego z czytelnikiem; adres urządzenia jest jednym z faktów, które ta
linia ma nieść. Adres bierze się z pola żądania HTTP, nie z nagłówków przekazywania, ponieważ te podaje
strona trzecia i można je napisać dowolnie, a dziennik ma nieść fakt gniazda, nie deklarację nadawcy.
Brak adresu daje wpis oznaczony jako nieustalony zamiast pustego miejsca, żeby czytający widział różnicę
między brakiem wiedzy a pominięciem.

## budowa/server/internal/protocol/zadanie.go
Okno komunikacji jako zasięg wygrywa z projektem, środowiskiem i pozostałymi
poziomami. Zasieg jest soczewką na ładunek, nie komunikatem: kontrakt nie ma
jednej struktury zasięgu, więc cztery poziomy wyjmuje się z ładunku po
nazwach pól, których kontrakt używa w treściach komend.

Klient koreluje odpowiedź po identyfikatorze żądania i po stanie: koperta bez
pola stanu nie rozstrzyga obietnicy wywołania, więc okno stoi w wiecznym
ładowaniu. Odmowa zeStanemOdmowy jest zwykłą odpowiedzią błędną; połączenie
nie jest zrywane, sesja nie jest blokowana, kolejne żądania są przyjmowane.
## budowa/server/internal/dane/studio.go
Wersje dokumentu i propozycje zmiany leżą w plikach sąsiednich tego samego repozytorium, rozdzielonych wedle odpowiedzialności; interfejs deklaruje wyłącznie ten plik, w całości, także metody obszarów wersji i propozycji, żeby cały kontrakt obszaru stał w jednym miejscu. Pole treści niesie treść krótką wprost, a pole odwołania niesie odwołanie do pliku dla treści obszernej odczytanej do tekstu, tym samym sposobem co w repozytorium wiadomości. Postać dokumentu, czyli drzewo postaci, arkusz stylów, sekcje, obiekty osadzone, aparat i pola, wchodzi zagnieżdżonym interfejsem osobnego repozytorium postaci, a nie przepisaniem jego metod: obszar postaci ma kilkadziesiąt metod i wypisanie ich po raz drugi znaczyłoby dwa wykazy jednego kontraktu, z których jeden prędzej czy później zostałby w tyle. Studio ma jedno repozytorium, nie dwa, więc adapter modułu dostaje postać tą samą zależnością, którą dostaje dokument.

## budowa/server/internal/dane/centrum_powiadomien.go

Centrum powiadomień jest trwałym rejestrem tych samych zdarzeń, których ulotną postacią jest
powiadomienie typu Toast. Rozgłoszenie zdarzenia i przekład na kształt kontraktu należą do
rdzenia, nie do tego repozytorium.

Odczytaj bez wskazanych identyfikatorów jest działaniem zbiorczym centrum: bierze wszystkie
zdarzenia nowe naraz.

Przywroc: bez tej operacji odłożenie zdarzenia na później byłoby cichym skasowaniem go
z rejestru.

## budowa/server/internal/transport/petla_odbioru.go

Kontekst rdzenia jest kontekstem serwera, nie połączenia: rozłączenie klienta nie przerywa pracy już
rozpoczętej przez rdzeń, sesja, okno i proces biegną dalej, a wynik trafi do pozostałych urządzeń konta
rozgłoszeniem. Źródło rdzenia jest dostawcą bieżącej realizacji obsługi komend, a nie wartością
zamrożoną w chwili nawiązania: rdzeń podłącza się po utworzeniu serwera, a połączenie nawiązane, zanim
to nastąpi, musiałoby wtedy trzymać pusty rdzeń na całe swoje życie. Pobranie rdzenia dopiero przy
obsłudze komunikatu sprawia, że komendy z takiego połączenia zaczynają być obsługiwane, gdy tylko
rdzeń zostanie podłączony. Straż jest rozstrzygnięciem o wystawieniu nasłuchu ustalonym raz, przy
normalizacji ustawień, i idzie przez pętlę odbioru do wykonania, ponieważ tam stoi jedyne wejście
żądania do rdzenia.

Odczyt w pętli odbioru kończy się na cztery sposoby: przekroczenie limitu odczytu zamyka gniazdo od
strony rdzenia, zatrzymanie rdzenia zamyka je z woli procesu, zerwanie sieci nie jest niczyją decyzją,
a odejście urządzenia zamyka je od jego strony. Linia dziennika jest jedynym trwałym śladem po
rozłączeniu, więc niesie to, co naprawdę zaszło, razem ze zdaniem biblioteki jako szczegółem.

Odczyt tożsamości z powitania dzieje się w tym pliku, a nie w rdzeniu, ponieważ tożsamość jest
własnością połączenia i mieszka przy nim; rdzeń ją czyta, a nie zapisuje. Gdyby zapisywał, musiałby
dostać ujście do ręki, a ujście do rdzenia świadomie nie idzie tą drogą. Odczyt jest czysty: transport
bierze pole, którego kształt i tak zna z kontraktu, i nie rozstrzyga o nim niczego. Odczyt dzieje się
przed oddaniem żądania rdzeniowi, więc zdarzenia rozgłoszone przez samo powitanie znają już klienta.
Ładunek nieczytelny albo pole puste zostawia tożsamość nietkniętą, ponieważ powitanie ma się udać
zawsze.

Każde żądanie w metodzie przyjmij idzie osobnym biegiem, więc komenda długotrwała nie zatrzymuje
odczytu, a komenda przerywająca dociera w trakcie jej wykonania. Odpowiedź niesie identyfikator
żądania, więc kolejność odpowiedzi nie ma znaczenia dla korelacji.

## budowa/server/internal/repozytorium/pliki.go
Bez sprowadzenia wskazania do katalogu roboczego, wejście w rodzaju odniesień
do katalogu nadrzędnego w nazwie pliku byłoby drogą do dowolnego miejsca na
dysku serwera, a okno modułu Developer stoi po to, żeby pracować w jednym
repozytorium. Sprawdzenie idzie po ścieżce rozwiniętej z dowiązań, nie po
samym napisie: dowiązanie symboliczne wskazujące poza katalog jest tym samym
wyjściem, tylko zapisanym inaczej.

Operator, który podaje ścieżkę zagnieżdżoną w Zaloz, prosi o plik pod tą
ścieżką, a nie o komunikat, że katalog pośredni nie istnieje; plik istniejący
nie jest nadpisywany, bo zakładanie nie jest zapisem i nie ma prawa skasować
cudzej treści.

Węzeł, którego nie było, nie zatrzymuje czynności Usun: Operator zaznaczający
kilkanaście plików nie ma dostawać odmowy przez jeden, który zniknął
wcześniej, ale nie ma go też w wykazie, bo tego pliku ta czynność nie
usunęła.

Czynności plikowe działają dzięki KorzenRoboczy także tam, gdzie Operator
jeszcze nie założył repozytorium.

## budowa/server/internal/dane/design_adnotacje.go

Adnotacja przeżywa zapis układu warstw celowo: pole adnotacji przy warstwie w kompozycji
wraca przepisane od nowa przy każdej aktualizacji układu, więc uwaga zostawiona przez jedną
osobę znikałaby przy pierwszym przesunięciu warstwy przez drugą, gdyby nie leżała osobno.

WarstwaID i NadrzednaID są identyfikatorami zewnętrznymi, nie kluczami wierszy, bo warstwa
bywa już usunięta z kompozycji, a odpowiedź w wątku zakłada się w jednym przebiegu okna
razem z adnotacją nadrzędną.

Autor adnotacji nie wchodzi w nadpisanie: adnotację zakłada jedna osoba, a zmiana treści
przez drugą nie czyni jej autorką cudzej uwagi.
## budowa/server/internal/dane/studio_adnotacje.go
Komentarz i adnotacja dzielą jedną tabelę, bo dzielą wszystkie kolumny poza jedną: komentarz wisi przy zakresie znaków treści, adnotacja przy numerze fragmentu porównania. Rozdzielenie ich na dwie tabele dałoby sześć kolumn powtórzonych i zmusiłoby pytanie o to, co ktoś napisał przy tym dokumencie, do sumy dwóch zapytań; rodzaj wiersza rozstrzyga osobna kolumna rodzaju. Warunek na stan oczekujący w poleceniu rozstrzygnięcia zmiany śledzonej jest zamierzony: decyzja raz podjęta nie zmienia się drugim wywołaniem, więc powtórzone przyjęcie tej samej zmiany oddaje wynik ujemny zamiast cicho nadpisywać odrzucenie przyjęciem.

## budowa/server/internal/dane/design_druk.go

Pola opcjonalne profilu druku niosą wskaźnik celowo: brak spadu znaczy wzięcie domyślnego
rdzenia, a spad zerowy znaczy druk bez spadu — to dwa różne rozstrzygnięcia i dwa różne
wyniki w drukarni.

## budowa/server/internal/transport/polaczenie.go

Zapis do gniazda prowadzi jedna pętla wysyłki, odczyt druga pętla odbioru; dzięki temu wysyłka z wielu
miejsc rdzenia równocześnie nie wymaga blokady na gnieździe i nigdy nie miesza ramek. Rdzeń widzi
połączenie wyłącznie przez interfejs Ujscie, czyli identyfikator, konto i wysyłkę koperty.

Chwila bieżąca jest w opisie metody Tozsamosc istotna: identyfikator klienta dochodzi dopiero
z powitaniem, więc żądanie wcześniejsze widzi tożsamość uboższą. To jest prawda o stanie, a nie brak do
naprawienia — rdzeń zamilknie wtedy o sprawcy zdarzenia zamiast go zmyślić.

## budowa/server/internal/repozytorium/roznica.go
Wynik różnicy wychodzi fragmentami zgodnymi z kontraktem, nie tekstem unified
diff, bo klient koloruje wiersze i przypina do nich uwagi przeglądu, więc
potrzebuje wierszy rozpoznanych po rodzaju, a nie napisu z myślnikiem na
początku.

Porównanie dwóch odwołań w paryPorownania nie sięga po katalog roboczy, bo
obie strony pochodzą z magazynu obiektów repozytorium.

Pliki niezmienione odpadają w paryDwochOdwolan przy zestawianiu drzew, a nie
przy liczeniu różnicy, inaczej odpowiedź niosłaby tysiąc pustych wpisów dla
repozytorium, w którym zmienił się jeden plik.

Licznik wszystkich zatwierdzeń w Historia rośnie także po osiągnięciu
granicy wyniku, bo odpowiedź niesie, ile zatwierdzeń jest w sumie, nie ile
pokazano; bez tego klient nie ma jak napisać, że pokazano część z całości.

## budowa/server/internal/dane/design_fotografia.go

Łańcuch edycji nie dubluje wariantów zasobu: że wariant powstał ze źródła, mówi już kolumna
wariantu zasobu; odczyt łańcucha idzie po zasobie wynikowym, więc każdy krok da się nazwać
i powtórzyć.

Czynność zawsze wstawia nowy wiersz: łańcuch edycji jest historią, a historia się nie
nadpisuje — ta sama czynność puszczona dwa razy na tym samym zasobie to dwa osobne ogniwa.

Zapytanie o czynność źródłową zasobu oddaje wykaz, choć wiersz jest najwyżej jeden: ten sam
zasób mógłby teoretycznie powstać dwiema drogami, a odczyt ma pokazać stan bazy, nie
założenie o niej.
## budowa/server/internal/dane/studio_cyfryzacja.go
Pozycja kolejki istnieje, zanim jakikolwiek dokument z niej powstanie, i bywa odrzucona, zanim taki dokument powstanie; wiązanie jej z dokumentem wymagałoby zakładania dokumentu pustego przy każdym wskazaniu pliku, także tym, które skończy się odmową rozpoznania. Warstwa słów rozpoznanych i bloki układu stoją tekstem w formacie JSON, bo nie są bytem samodzielnym: nie mają własnego cyklu życia, nikt się do nich nie odwołuje z zewnątrz i giną razem z pozycją, więc tabela podrzędna dałaby wyłącznie złączenie przy każdym odczycie bez żadnej korzyści w zamian.

## budowa/server/internal/transport/rejestr_polaczen.go

Jedno konto ma wiele urządzeń równocześnie, więc rejestr nie jest odwzorowaniem konta na połączenie,
lecz zbiorem połączeń przeszukiwanym po koncie. Konto połączenia zmienia się w czasie jego życia, więc
rejestr nie kopiuje go do własnego indeksu, tylko odczytuje wprost z połączenia przy każdym zapytaniu.
Metoda wszystkie zwraca kopię zbioru połączeń, nie samo odwzorowanie, ponieważ wysyłka idzie poza
blokadą, więc rozłączenie w trakcie rozgłoszenia niczego nie zakleszcza.

## budowa/server/internal/dane/design_gradienty.go

Wskazanie ścieżki albo warstwy pominięte zapisuje się pustym napisem, nie wartością pustą
bazy: silnik bazy liczy dwie wartości puste za różne w indeksie unikalnym, więc gradient
całej kompozycji zakładany dwa razy powstałby dwa razy zamiast nadpisać się raz.

Identyfikator zewnętrzny gradientu nadaje wołający i służy wyłącznie temu, żeby gradient
nowy dostał własną nazwę; rozstrzyga cel zapisu, nie identyfikator — gradient wysłany drugi
raz na tę samą warstwę nadpisuje ten, który tam stoi, choćby przyszedł z nowym
identyfikatorem.
## budowa/server/internal/dane/studio_praca_widok.go
Repozytorium sąsiednie zna ten sam wiersz od strony autozapisu: odstęp, zdarzenia okna, wygasanie kopii, skutek ostatniego zapisu. Kolumny widoku, czyli tryb powierzchni, skala, linijki, układ stron, przewijanie, podświetlenie zmian wykonawcy i przybornik, pytane są przez zupełnie inną parę komend, nigdy razem z nastawami autozapisu. Osobny odczyt tych samych wierszy nie zakłada drugiego pojęcia nastawy: tabela jest jedna, wiersz zakłada się jedną drogą, a każda z dwóch grup kolumn ma własne polecenie zapisu, bo jedno polecenie na obie grupy kazałoby widokowi przepisywać nastawy autozapisu, których nie zmieniał. Wołający sięga najpierw po zapis pełnej nastawy pracy, potem po ten odczyt widoku, bo dwie drogi zakładania tego samego wiersza rozjechałyby się przy pierwszej zmianie wartości domyślnej.

## budowa/server/internal/repozytorium/szukanie.go
Wyszukiwanie globalne jest czynnością, bez której Code Editor przestaje być
edytorem kodu, więc silnik regexp nie może zależeć od zewnętrznego programu,
którego instalka nie niesie. RE2 nie ma odwołań wstecznych, i to jest jedyna
różnica, którą Operator zobaczy: wyrażenie z odwołaniem wstecznym dostanie
odmowę nazwaną, nie ciche zero trafień. Wyszukiwanie, które wchodzi w katalog
zależności i katalog budowania, oddaje tysiące trafień w kodzie, którego
Operator nie pisał, i jest wtedy bezużyteczne, choć formalnie poprawne.

Granica domyślna istnieje, bo wyszukanie pojedynczej litery w dużym
repozytorium dałoby zbiór, którego klient nie postawi na ekranie, a rdzeń
trzymałby go w pamięci w całości.

Trzecia wartość zwracana przez Szukaj mówi, czy wynik przycięto granicą;
druga niesie, ile trafień znaleziono przed przycięciem — sam wykaz przycięty
bez tej liczby kazałby Operatorowi zgadywać, czy zobaczył wszystko.

Podgląd jest stanem domyślnym pracy eksperckiej w Zamien: zamiana masowa
dotyka wielu plików naraz i cofnięcie po niej nie istnieje. Dlatego czynność
oddaje wykaz zmian razem z odpowiedzią o tym, czy je zapisała, a nie samo
słowo wykonano, z którego Operator nie wyczyta, co się stało z jego kodem.

Katalog bez repozytorium też podlega przeszukaniu w obszarSzukania: Operator
otwiera w oknie także katalogi, których nie wersjonuje, a odmowa wyszukiwania
byłaby wtedy odmową bez powodu; reguły pomijania są wtedy puste poza
katalogiem repozytorium.

Plik ignorowanych ścieżek korzenia niesie reguły katalogów budowania
i zależności, czyli te, których pominięcie decyduje o użyteczności wyniku.
Reguły podkatalogów zawężają wynik dodatkowo, a ich brak oznacza wyłącznie
kilka trafień więcej, nigdy mniej.

Dopasowanie wzorca w pasujeGlob idzie dwutorowo, bo Operator pisze wzorce
zarówno rozszerzenia plików, jak i ścieżek katalogów; dopasowanie ścieżek nie
zna podwójnej gwiazdki, więc taki wzorzec sprowadza się do przedrostka
ścieżki — to jest znaczenie, którego Operator się spodziewa.

PlikiDoPrzejrzenia jest wystawione poza pakiet dla skanowania bezpieczeństwa
i jakości: skan przechodzi dokładnie ten sam zbiór plików, co wyszukiwanie,
więc jedno przejście ma stanowić jedną prawdę o tym, co należy do
repozytorium. Skan zgłaszający sekret w katalogu pobranych zależności byłby
wykazem, którego nikt nie czyta. Wzorce włączające zawężają wynik tak samo,
jak w Szukaj; pusty wykaz znaczy całe drzewo.

## budowa/server/internal/dane/design_ikony.go

Katalog ikon otwartoźródłowych bazy nie dotyka: jest wkompilowany w binarium rdzenia, więc
rdzeń zna go bez kroku zasiewu, a stanowisko bez ani jednej ikony własnej i tak ma czego
szukać.

Nazwa ikony jest jedyna w oknie — zderzenie wychodzi jako błąd bazy, nie jako cicha podmiana
cudzej ikony: dwie ikony o jednej nazwie w jednym oknie dałyby sprite z dwoma symbolami
o tym samym identyfikatorze, czyli plik, którego przeglądarka nie złoży.

Etykiety ikon idą po zamknięciu kursora zapytania, tak samo jak przy kolekcjach.

## budowa/server/internal/transport/serwer.go

Rdzeń podłącza się po utworzeniu serwera, a nie przy jego budowie, żeby zerwać zależność cykliczną:
rdzeń potrzebuje rozgłośni transportu, transport potrzebuje rdzenia do obsługi komend, a żaden pakiet
nie importuje pakietu drugiego. Reguła bez blokad domyślnych mówi, że brak nastawy nie wstrzymuje
pracy, ale niekompletna para plików TLS nastawą nie jest, tylko połową wskazania — praca otwartym
tekstem przy wskazanym certyfikacie byłaby cichym zejściem poniżej tego, o co poprosił operator maszyny.
## budowa/server/internal/dane/studio_praca_zmiany.go
Repozytorium sąsiedni zna zmianę śledzoną sprzed dobudowy: rodzaj, autora grubym rozróżnieniem człowiek-model, zakres, brzmienie przed i po, decyzję. Kolumny tożsamości agenta, podagenta oraz postaci przed i po zmianie dołożyła późniejsza migracja i pyta o nie wyłącznie ten odcinek, bo przełącznik pokazujący wszystko, co zrobił model, musi rozdzielić dwóch agentów pracujących naraz, a nie pokazać obu jako jednego. Dopisanie ich do pliku sąsiedniego byłoby wejściem w plik cudzego odcinka; osobny odczyt tych samych wierszy nie zakłada drugiego pojęcia zmiany śledzonej, bo tabela jest jedna, wiersz zakłada zapis zmiany śledzonej, a te kolumny stempluje się na wierszu już istniejącym.

## budowa/server/internal/transport/statyka.go

Katalog niewskazany albo jeszcze niezbudowany nie wstrzymuje nasłuchu ani kanału WebSocket — żądanie
pliku dostaje wtedy odpowiedź o braku pliku z wyjaśnieniem, a rdzeń pracuje dalej. Katalog sprawdzany
jest przy każdym żądaniu, więc zbudowanie klienta po starcie rdzenia wystarcza, by pliki zaczęły się
serwować bez ponownego uruchomienia.

## budowa/server/internal/session/bledy.go
Odwzorowanie błędów pakietu na kody kontraktu leży w tym pakiecie: pakiet
session nie zakłada własnego katalogu kodów.
## budowa/server/internal/dane/studio_propozycje.go
Propozycja nie jest wersją: porównanie różnic przyjmuje identyfikator propozycji zamiennie z identyfikatorem wersji docelowej, więc porównanie czyta propozycję i wersję tym samym mechanizmem odczytu treści, ale zapis obu bytów jest rozdzielony, bo zatwierdzenie propozycji do repozytorium wykonuje osobny zapis dokumentu, nie ten plik. Fragmentów różnicy tu nie ma: liczą się w locie z dwóch treści w warstwie rdzenia, a ten plik oddaje wyłącznie treść propozycji do porównania, nie sam wynik porównania.

## budowa/server/internal/dane/design_kolekcje.go

Liczba zmienionych przypisań kolekcji liczy się z wyniku poleceń, nie z długości nadesłanego
wykazu: kontrakt pyta, ile przypisań naprawdę się zmieniło, a dołożenie zasobu już
należącego do kolekcji zmienia zero.

Dołożenie zasobu, który już jest w kolekcji, zmienia zero przypisań i tak ma zostać
policzone — powtórzenie w wykazie nie jest odmową.

Znacznik czasu kolekcji przestawia się tylko wtedy, gdy coś się naprawdę zmieniło: wykaz
kolekcji sortuje się po nim, a przypisanie bez skutku nie ma prawa przestawiać kolejności
na ekranie.

## budowa/server/internal/session/identyfikator.go
Jedna sesja prowadzi wiele okien i wiele procesów biegnących równolegle,
każde z własnym kanałem modelu i własnym katalogiem roboczym.

## budowa/server/internal/dane/design_komponenty.go

Instancja komponentu wskazuje warstwę identyfikatorem zewnętrznym, który przeżywa zapis
planszy. Liczba instancji przy komponencie jest liczona w bazie przy odczycie, nie
przechowywana obok jako osobna wartość: licznik trzymany osobno rozjeżdżałby się
z rzeczywistością przy pierwszym usunięciu wiersza, o którym nikt licznika nie powiadomił.

## budowa/server/internal/session/izolacja.go
Rozstrzyganie zasięgu w pakiecie konfig mówi, jaka wartość izolacji obowiązuje
w danym oknie, ale samo z siebie nie zmienia niczego w wykonaniu; egzekutor
bierze tę wartość i odrzuca uruchomienie, które by ją naruszyło, a wartość
wyłączona nie ogranicza niczego, bo stanem wyjściowym platformy jest pełna
swoboda operacyjna.

Egzekutor nie zamyka procesu modelu w piaskownicy systemu operacyjnego. Panuje
wyłącznie nad tym, co procesowi podaje: katalogiem startowym, środowiskiem,
katalogiem danych kanału, wykazem korzeni, adresami i treścią tury. Gniazd
sieciowych ani wywołań systemowych uruchomionego już procesu potomnego nie
ogranicza, do tego potrzebne są środki systemu operacyjnego. Nazwy punktów
izolacji pochodzą z pakietu konfig.

Porównanie ścieżek w SciezkaWewnatrz idzie po członach ścieżki, nie po
prefiksie napisu: katalog nazwany podobnie z dopiskiem nie leży w katalogu
oryginalnym. Na systemie plików nierozróżniającym wielkości liter porównanie
też jej nie rozróżnia, inaczej ta sama ścieżka zapisana inną wielkością
omijałaby obszar.
## budowa/server/internal/dane/studio_strona_znaki.go
Nie ma tu tablicy znaków: nazwy znaków, ich punkty kodowe i grupy są wiedzą rdzenia, tak samo jak arkusz stylów fabryczny i wykaz nośników druku, więc do bazy schodzi wyłącznie to, co Operator zmienił albo czym się posłużył. Drugi wykaz znaków w tabeli rozjechałby się z wykazem rdzenia przy pierwszym uzupełnieniu. Wykaz zasad autozamiany jest odwrotnie: stoi w bazie, także w części fabrycznej, bo migracja tak go założyła i bo Operator ma prawo zasadę fabryczną wyłączyć; powtórzenie wykazu fabrycznego w kodzie rdzenia dałoby dwie prawdy o tym, co wchodzi w miejsce skrótu. Zmiana zasady fabrycznej nie zdejmuje jej oznaczenia: Operator, który wyłączył zasadę fabryczną, nadal ma przed sobą zasadę fabryczną wyłączoną, a nie zasadę własną, którą wolno usunąć; bez tego rozróżnienia dałoby się usunąć zasadę fabryczną, wpisując ją najpierw jako własną, czyli obejściem. Nazwy pomocnicze tego pliku niosą wspólny przedrostek, bo przestrzeń nazw pakietu jest dzielona z innymi wykonawcami.

## budowa/server/internal/transport/tozsamosc.go

Tożsamość mieszka przy połączeniu, tam, gdzie już mieszka konto, i wchodzi do kontekstu żądania tą samą
jedną drogą, którą wchodził identyfikator gniazda. Nie powstaje ani drugi wpis kontekstu, ani drugi
rejestr, ani pole w kopercie kontraktu, ponieważ koperta jest zamrożona, a tożsamość i tak nie jest
własnością komunikatu, tylko własnością łącza. Fakty biorą się z dwóch źródeł składanych w jeden skład:
przy nawiązaniu serwer narzędzi modelu przedstawia się parametrami zapytania przy zestawianiu gniazda,
tą samą drogą, którą urządzenie od zawsze wskazuje konto, ponieważ jest klientem wołającym komendy,
a nie oknem interfejsu, więc powitania nie wysyła; przy powitaniu okno interfejsu niesie identyfikator
klienta w komunikacie powitalnym, a transport odczytuje je z ładunku powitania i dokłada do tożsamości
połączenia. Transport niczego nie rozstrzyga: nie zna pojęcia rodzaju sprawcy, nie wie, co to operator
ani asystent, i nie ma w tym pliku ani jednej wartości wyliczenia kontraktu — niesie fakty, a
rozstrzygnięcie, czyja to ręka, należy do rdzenia, bo tylko rdzeń zna rolę okna. To nie jest uprawnienie:
tożsamość nie rozstrzyga ani razu, czy coś wolno, tym zajmuje się wyłącznie straż bramki i pyta
o zupełnie co innego; parametr podany przez wołającego jest tu opisem, więc jego podrobienie niczego
nie otwiera — gdyby cokolwiek od niego zależało, byłby bramką.

Parametr klienta nazywa parametr zapytania i nagłówek, którymi klient przedstawia swój identyfikator
już przy nawiązaniu; okno interfejsu go nie używa, jemu wystarczy powitanie, ale klient bez powitania,
czyli serwer narzędzi modelu, innej drogi nie ma. Parametr rodzaju niesie rodzaj klienta, czym jest
program po drugiej stronie gniazda, a wartość znaną transportowi jest jedna — rodzaj narzędzi. Parametr
zasięgu niesie rolę okna, w którego imieniu pracuje serwer narzędzi; transport wartości tej nie zna
i nie porównuje, przenosi napis, bo import w tę stronę odwróciłby zależność. Parametr okna niesie okno
rozmowy serwera narzędzi — do rozstrzygnięcia sprawcy niepotrzebne, do zrozumienia dziennika konieczne.
Stała oznaczająca rodzaj narzędzi stoi w tym pliku, a nie w pakiecie narzędzi, żeby obie strony rozmowy,
ta, która napis wysyła, i ta, która go czyta, brały go z jednego miejsca.

Rdzeń ma milczenie pola tożsamości przenieść dalej — kontrakt mówi o polu sprawcy wprost, że brak
znaczy, iż rdzeń nie potrafił tego rozstrzygnąć; „nie wiadomo" to stan zwykły dla gniazda, które jeszcze
się nie przedstawiło.

Parametr zapytania w funkcji tozsamoscZadania stoi przed nagłówkiem, tak samo jak przy koncie, ponieważ
klient WebSocket w przeglądarce nagłówków ustawić nie potrafi, a klient biblioteczny potrafi obu dróg.
Brak wskazania nie odrzuca nawiązania: gniazdo bez tożsamości pracuje dalej, tyle że jego zdarzenia
pójdą bez sprawcy.

Funkcja ParametryTozsamosci stoi w tym pliku, przy odczycie, a nie po stronie klienta, ponieważ napis
parametru ma mieć jedno źródło dla piszącego i czytającego; klient podaje wartości, nazw nie zna
i nie powtarza. Wartości puste nie wchodzą do wyniku — parametr pusty i nieobecny mają znaczyć to samo.

## budowa/server/internal/session/izolacja_polecenie.go
Miejsce egzekucji izolacji jest jedno: polecenie zbudowane przez warstwę
kanału, zanim pójdzie do uruchamiacza. Polecenie niosące ścieżkę albo
środowisko spoza obszaru okna zostaje odrzucone, proces z takim poleceniem
nie startuje.

Wartość pusta, którą izolacja jednoznacznie wyznacza, zostaje uzupełniona
w SprawdzPolecenie: brak wskazania nie jest naruszeniem, bo brak danych ma
dawać poprawny wynik, nie awarię. Wskazanie wychodzące poza obszar okna jest
naruszeniem, bo jest decyzją sprzeczną z ustawieniem Operatora.

SprawdzKatalogDanych obsługuje zarówno sprawdzenie polecenia, jak i punkt,
w którym kanał wybiera katalog konfiguracji konta — reguła jest jedna.
Wskazanie puste jest naruszeniem: kanał bez wskazania sięga po katalog
wspólny, a właśnie tego izolacja zabrania.

## budowa/server/internal/session/izolacja_przydzial.go
Każdy z trzech zakresów izolacji ma tę samą treść: włączony znaczy zasób
dedykowany oknu, wyłączony znaczy zasób wspólny platformy. Egzekucja jest
jedna dla trzech zakresów: zasób, którego właścicielem nie jest to okno,
zostaje odrzucony. Właściciela wskazuje ten, kto zasób przydziela — rejestr
procesów kluczowany oknem, pula kont oddająca kod profilu, przydział serwera
wykonania.
## budowa/server/internal/dane/studio_wejscie_pochodzenie.go
Wiersz pochodzenia zakłada ten, kto fragment wnosi: wniesienie pliku, obrazu, fragmentu z biblioteki albo ze strony sieci. Wykaz pochodzenia czyta go inny odcinek repozytorium, dlatego odczyt jest tu równie pełny jak zapis, żeby tamten odcinek nie musiał zakładać drugiego. Zakres jest liczony w znakach, tak samo jak zakres blokady; fragment usunięty zostawia wiersz o zerowej długości świadomie, bo to, że Operator wniósł kiedyś fragment z danego źródła, jest faktem, którego usunięcie tekstu nie unieważnia.

## budowa/server/internal/dane/design_kompozycje.go

Zapis kompozycji nadsyła zawsze całą listę warstw naraz — kontrakt nie ma trybu częściowej
zmiany — dlatego ZapiszKompozycje usuwa warstwy zastane i wstawia przysłany komplet od nowa,
inaczej usunięcie warstwy w oknie nie usunęłoby jej w bazie i kompozycja rozeszłaby się
z tym, co widać na ekranie.

Brak identyfikatora zewnętrznego przy zapisie zakłada kompozycję nową jedną ścieżką zapisu:
zapis idzie po identyfikatorze z klauzulą podmiany przy konflikcie, a brak identyfikatora
rozstrzyga wywołujący, generując nowy przed wywołaniem tej metody — repozytorium zna
wyłącznie tryb załóż-albo-nadpisz.

Kolejność warstw rozstrzyga porządek listy, bo znacznik czasu ma rozdzielczość milisekundy
i dwa zapisy z jednej pętli okna potrafią w nią trafić razem.

Kompozycje zwraca wykaz bez granicy strony: kontrakt wykazu kompozycji nie niesie ani
limitu, ani liczby całkowitej liczonej osobno od długości wykazu — liczba odpowiedzi jest
liczbą zwróconych kompozycji, więc przycięcie wykazu tutaj rozjechałoby obie liczby naraz.
Okno puste oddaje wykaz pusty, nie wszystkie kompozycje: brak wskazania okna sprawdza
wołający, a zapytanie i tak porównuje kolumnę z pustym tekstem, którego żadne okno nie nosi.

## budowa/server/internal/transport/ustawienia.go

Pojemność kolejki domyślna to liczba komunikatów oczekujących na zapis do jednego gniazda; bufor chroni
rdzeń przed zablokowaniem na wolnym urządzeniu, a przepełnienie kończy pojedynczą wysyłkę, nie całą
sesję. Konto domyślne obowiązuje, dopóki urządzenie nie wskaże konta, ponieważ uwierzytelnianie jest
jedyną kontrolą dostępu i w fazie budowy nie działa, więc brak konta nie może wstrzymać połączenia.
Nagłówek konta jest nagłówkową postacią parametru konta, dla klientów, które nie mogą dopisać parametru
do adresu. Wyjście poza pętlę zwrotną jest osiągalne jednym polem i nadal ostrzega w dzienniku; stała
adresu domyślnego jest wewnętrzna, ponieważ wołający wskazują adres wprost, nie sięgają po wartość
domyślną.

Pole Adres rdzenia zmienia umiejscowienie w czasie, więc droga do wystawienia zostaje otwarta —
zmienia się wyłącznie to, co dzieje się bez jawnego wskazania. Pole WszystkieInterfejsy jest osobnym
polem, a nie pustym napisem, ponieważ brak wskazania i chęć wystawienia wszędzie to dwa różne zdania,
mające wyglądać różnie w miejscu wywołania. Pusty wykaz pochodzeń dozwolonych bierze pochodzenia
własne opisane w warstwie nawiązania połączenia.

Pole WymogLogowania ma trzy stany, nie dwa, dlatego jest wskaźnikiem: brak wskazania oddaje
rozstrzygnięcie adresowi nasłuchu, czyli bez wymogu na pętli zwrotnej i z wymogiem przy nasłuchu
szerszym; wartość prawda włącza wymóg także na pętli zwrotnej; wartość fałsz znosi wymóg także przy
nasłuchu szerszym, co jest dozwolone, ale nigdy ciche, bo dziennik mówi o tym wprost, kiedy maszyny
operatora stoją otworem przez tor zdalny. Wartość logiczna zamiast wskaźnika kasowałaby różnicę między
brakiem wskazania a jawnym wskazaniem odmowy, a to jest tu cała różnica.

Brak obu plików pary TLS zostawia otwarty tekst i, poza pętlą zwrotną, ostrzeżenie w dzienniku.
Wskazanie tylko jednego z dwóch plików jest błędem konfiguracji i zatrzymuje start, ponieważ cicha
praca otwartym tekstem przy wskazanym certyfikacie byłaby najgorszym z możliwych wyników. Wartość
ujemna portu nie występuje, bo konfiguracja sprawdza jej zakres wcześniej.

Wpięcie rozpoznania wystawienia przy samym otwarciu gniazda sieciowego byłoby bliżej faktu, ale metoda
adresNasluchu wołana jest kilka razy i ostrzeżenie by się dublowało; ostrzeżenie idzie po ustaleniu
dziennika, żeby brak dziennika kierował je do kosza, a nie gubił wywołania. Kto chce wystawienia
szerszego, mówi to wprost jednym z dwóch pól ustawień.

## budowa/server/internal/transport/wystawienie.go

Cztery decyzje rdzenia były rozsądne osobno, dopóki nasłuch nie wychodził poza pętlę zwrotną: nasłuch
pusty oznaczał wszystkie interfejsy, uwierzytelnianie celowo nie działało, wzorce pochodzenia otwarte
nie sprawdzały niczego, a warstwa TLS nie istniała. Trzy z nich są już zdjęte: adres pusty znaczy pętlę
zwrotną, pochodzenie jest sprawdzane wykazem wzorców, a warstwa TLS włącza się parą plików wskazaną
przełącznikiem albo zmienną środowiska. Czwartą, brak uwierzytelniania, wygasza straż bramki: poza
pętlą zwrotną gniazdo wykonuje wyłącznie powitanie, logowanie i rejestrację, dopóki nie przedstawi
tokenu sesji, a na pętli zwrotnej straż nie istnieje. Ostrzeżenie zostaje, ale mówi o stanie bieżącym:
składa się z tego, co naprawdę zastane w nastawach, nie z braków, których już nie ma. Ten plik nie
zatrzymuje startu — rdzeń ma wstać i działać, tylko ma powiedzieć głośno, na czym staje.

Linia ostrzeżenia wystawienia mówi stan faktyczny: co od tej chwili obowiązuje i co nadal zostaje na
operatorze maszyny. Nie jest to już ostrzeżenie o braku bramki, bo bramka jest, lecz zawiadomienie
o zmianie zachowania rdzenia — nic istotnego nie dzieje się po cichu.

Bez warstwy szyfrowanej zdanie ostrzeżenia jest mocniejsze: token sesji jedzie tym łączem przy każdym
powitaniu, a otwartym tekstem jedzie w postaci czytelnej dla każdego po drodze. Warstwa szyfrowana nie
rozstrzyga, kto się łączy, od tego jest bramka, ale bez niej bramka broni wejścia, którego klucz leci
obok, na wierzchu.

Pusty adres nasłuchu to nie pętla zwrotna: otwarcie gniazda z pustym hostem wiąże wszystkie interfejsy,
zarówno IPv4, jak i IPv6, więc pusty adres jest najszerszym z możliwych wystawień, nie najwęższym. To
samo dotyczy zapisanych wprost adresów wszystkich interfejsów w obu wersjach protokołu. IPv4 i IPv6 idą
jedną drogą sprawdzania pętli zwrotnej, obejmującą całą sieć adresów lokalnych oraz ich zagnieżdżoną
postać w IPv6; identyfikator strefy adresu łączowego odcinamy przed rozbiorem, bo rozbiór go nie
przyjmuje, a bez odcięcia adres łączowy zostałby wzięty za nazwę. Nazwa, która nie jest adresem, jest
traktowana jako wystawienie, nie jako pętla zwrotna — jedynym wyjątkiem jest nazwa localhost i nazwy
w jej domenie, których rozwiązanie na pętlę zwrotną gwarantuje odpowiednia norma internetowa. Rdzeń nie
pyta o to systemu nazw, ponieważ odpytanie DNS przy starcie wstrzymywałoby start, a ostrzeżenie ma być
tanie i pewne; wynik z tego jest asymetryczny celowo — przy wątpliwości ostrzeżenie pada, bo ostrzeżenie
zbędne kosztuje linię dziennika, a ostrzeżenie pominięte kosztuje wystawiony rdzeń.

Funkcja ostrzezJezeliWystawiony nic nie zwraca i nic nie zatrzymuje — jedynym skutkiem jest linia
w dzienniku, ponieważ brak zabezpieczenia nie wstrzymuje startu i nic istotnego nie dzieje się po cichu.
Operator ma prawo zdjąć dźwignię wymogu logowania, bo to jego maszyna i jego rozstrzygnięcie, ale nie ma
prawa zrobić tego po cichu: nieuwierzytelnione połączenie sięga po hosty zdalne, czyli po zdalny serwer
i komputer w biurze, więc linia ostrzeżenia idzie zamiast zawiadomienia zwykłego, bo mówi o tym samym
nasłuchu rzecz ważniejszą.

## budowa/server/internal/zdalne/budzik_powiadomien.go

Bez tej pętli ponowienie i wygaśnięcie powiadomień byłyby trzema kolumnami, których nikt nigdy nie
rusza, czyli atrapą wymagania, a nie jego spełnieniem. Wzorzec jest wzięty z już działającej pętli
rdzenia: zegar taktujący, oczekiwanie na kontekst życia procesu, pierwszy przebieg od razu po starcie
i awaria jednego przebiegu, która nie zatrzymuje pętli. Pierwszy przebieg od razu ma znaczenie akurat
tutaj: powiadomienia zgłoszone tuż przed postojem rdzenia mają dolecieć zaraz po jego powrocie, a nie
po pełnym takcie.

Takt gęstszy niż trzydzieści sekund nie przyspieszyłby niczego, bo i tak czeka się na kolumnę kolejnej
próby, a takt rzadszy opóźniałby pierwsze podejście o więcej, niż wynosi cała jego zwłoka. Kolejność
w metodzie przebieg jest rozmyślna: wygaszanie idzie pierwsze, żeby przeterminowane powiadomienie nie
zdążyło polecieć w tym samym takcie, w którym straciło ważność, ponieważ takt to zawsze jakiś kawałek
czasu, a gdyby wysyłka szła pierwsza, budzik zabrzmiałby o sprawie, o której sam za chwilę orzeka, że
jest nieaktualna. Wygaszanie idzie także wtedy, gdy nadajnika nie ma, ponieważ sprawa nieaktualna jest
nieaktualna niezależnie od tego, czy było komu ją zanieść.

Brak nadajnika ma być widoczny w dzienniku, a nie odgadywany z tego, że nic nie dolatuje.

## budowa/server/internal/zdalne/hosty.go

Dwa odczyty bazy prowadzą do decyzji toru: nazwa hosta z ustawienia wykonania i wiersz hosta ze zgodą
w tabeli hostów zdalnych, a każda brakująca część drogi kończy się odmową trójczęściową.

Funkcja hostOkna odczytuje nazwę hosta wykonania po dwóch poziomach zasięgu, które pakiet zna
z tożsamości okna: zapis na samym oknie oraz poziom globalny, przy czym węższy poziom wygrywa. Poziomy
pośrednie, takie jak sesja, projekt czy moduł, zna wyłącznie rozstrzygacz rdzenia — pełne rozstrzyganie
wszystkich poziomów wymaga osobnego wpięcia, a do tego czasu odczyt węższego i najszerszego poziomu jest
uczciwym podzbiorem, nie atrapą całości.

Brak wiersza hosta i brak zgody na jego użycie są odmowami trójczęściowymi — każda nazywa dokładnie ten
ruch operatora maszyny, który ją zdejmuje.

## budowa/server/internal/zdalne/pliki.go

Każdy wykonany ruch bajtów zostawia wiersz prowenancji w tabeli przeniesień zdalnych. Nagrania dźwięku
nie jadą tym torem: dźwięk nie opuszcza maszyny operatora, a potrzeby też nie ma, ponieważ silnik mowy
jest usługą rdzenia i bierze ścieżkę na maszynie silnika, więc transkrypcja domyka się przed torem,
a do procesu zdalnego jedzie wyłącznie tekst. Wykaz rozszerzeń nagrań jest jeden, wspólny ze słownikiem
formatów nagrań silnika mowy — drugiej listy ten plik nie zakłada. W rdzeniu funkcję przenoszenia woła
spoina katalogów roboczych przy zasięgu zdalnym.

## budowa/server/internal/zdalne/polecenie.go

Tryb wsadowy wyklucza pytania interaktywne, ponieważ zgoda jest w bazie, nie w terminalu, a proces okna
nie ma przy sobie nikogo, kto by odpowiedział. Przyjęcie nowego klucza zapisuje klucz hosta przy
pierwszym połączeniu i odmawia przy jego zmianie, więc podmieniona maszyna nie dostanie procesu po
cichu. Limit czasu połączenia zamienia wieczne wiszenie na odmowę z podanym powodem.

Komenda zdalna ma postać zmiany katalogu i podmiany procesu powłoki zmiennymi środowiska, programem
i argumentami. Każdy człon jest cytowany zgodnie z regułami powłoki POSIX, więc treść polecenia nie
może zmienić kształtu komendy. Podmiana procesu powłoki oddaje procesowi miejsce powłoki — sygnał
zerwania połączenia trafia wprost do niego, nie do pośrednika.

## budowa/server/internal/zdalne/powiadomienia.go

Silnik doręcza po łączu, które operator już otworzył; wygaszonego telefonu nie budzi, ponieważ
wymagałoby to usługi wypychania powiadomień spoza tego systemu. Przy telefonie wygaszonym powiadomienie
czeka w kolejce i doleci przy najbliższym otwarciu; interfejs musi to pokazywać, bo inaczej operator
liczyłby na dzwonek, którego nie ma.

Format znacznika czasu nie używa skrótu obcinającego zera końcowe, ponieważ przy takim zapisie chwila
równa co do milisekundy wypadłaby raz z pełnymi zerami, jak zapisuje je baza, a raz bez nich, jak
zapisałby je język programowania, a terminy w tej kolejce porównuje się jako napisy. Porównanie
leksykalne tych dwóch zapisów dawałoby odpowiedź odwrotną do prawdziwej, więc próba wypadałaby
o milisekundę za wcześnie albo za późno bez żadnego śladu.

Po wyczerpaniu wykazu odstępów ponowienia obowiązuje ostatni odstęp: kolejka nie dobija urządzenia
częściej, ale też nie przestaje próbować przed terminem ważności. O tym, kiedy przestać, rozstrzyga
termin ważności powiadomienia, a nie licznik prób. Pole nadawania trzyma nadajnik podany przez
kompozycję i stoi osobno od zasilenia bazy, bo uchwyt bazy i droga doręczenia to dwa niezależne
wpięcia.

Funkcja Zglos wnosi powiadomienie do kolejki i oddaje jego klucz; sama droga nie doręcza, wysyłką
zajmuje się takt wywoływany funkcją Wyslij. Gdyby zgłoszenie doręczało od razu, decyzja podjęta w tej
samej chwili nie zdążyłaby powiadomienia odwołać.

Identyfikator kanału jest identyfikatorem klienta, tym samym napisem, który transport niesie jako
identyfikator klienta w tożsamości połączenia. Rejestracja powtórzona odświeża wiersz zamiast zakładać
drugi, więc powrót urządzenia na łącze może ją wołać bez sprawdzania, czy już było.

Sama droga odwołania nie wystarcza, by nie zawołać po fakcie: brak odstępu, w którym dałoby się coś
przegapić między odwołaniem a taktem, bierze się stąd, że takt sprawdza stan pod zamkiem zapisu na tym
samym wierszu. Bez tamtego zamka wywołanie odwołania byłoby wyścigiem.

Funkcja Wygas zamyka powiadomienia, którym minął termin ważności, i oddaje ich liczbę. Droga jest osobna
od funkcji Wyslij, ponieważ termin ważności ma mijać także wtedy, gdy nadajnika nie ma.

Brak nadajnika w funkcji Wyslij nie jest ciszą ani zerem: gdyby przebieg bez nadajnika po prostu nikomu
nie doręczył, każde powiadomienie podbijałoby licznik prób i po kilkunastu taktach wygasłoby, nigdy nie
mając odbiorcy, a kolejka twierdziłaby, że próbowała. Dlatego przebieg bez nadajnika odmawia całością
i nie tyka ani jednego wiersza: powiadomienia zostają z pełnym budżetem prób na chwilę wpięcia
nadajnika. Awaria taktu przerywa przebieg; podsumowanie oddaje to, co zdążyło się rozstrzygnąć przed
błędem.

Funkcja rozeslij woła każdy czynny adres i oddaje klucze rejestracji, które kopertę przyjęły; urządzenie
nieobecne odpowiada fałszem, co nie jest błędem przebiegu, tylko powodem ponowienia.

## budowa/server/internal/zdalne/tor.go

Droga rozstrzygnięcia toru ma pięć ogniw i każde brakujące jest osobną, nazwaną odmową: zasilenie bazy
rdzenia, host wskazany ustawieniem wykonania na poziomie okna albo globalnym, host wpisany do wykazu
hostów zdalnych, zgoda operatora na tym wierszu wydana oraz program SSH obecny na maszynie rdzenia.
Odmowa nie jest bramką wobec operatora maszyny: każda mówi, co się nie stało, dlaczego, i którym ruchem
operator to zmienia. Zgoda per host chroni maszyny operatora — rdzeń nie zainicjuje połączenia z maszyną,
której mu nie oddano.

## budowa/server/internal/dane/terminal_odczyt.go
Filtr dziennika procesów terminala idzie parametrem, nie sklejaniem tekstu SQL:
jedno przygotowane zapytanie obsługuje cztery zawężenia naraz, bo pusty
parametr znaczy „nie zawężaj”. Dzięki temu pamięć podręczna zapytań ma jedną
pozycję zamiast szesnastu, a wartości nigdy nie wchodzą do treści zapytania.

Uchwyt bazy w repozytoriumTerminala stoi obok przygotowanych zapytań, bo dwa
zapisy tego obszaru obejmują więcej niż jedno polecenie i muszą pójść jedną
transakcją: nadanie numeru wersji pozycji biblioteki wraz z wpisem tej wersji
oraz odpięcie klucza od wpisów hostów wraz z odczytaniem, których wpisów to
dotyczyło.

## budowa/server/internal/zdalne/zdalne.go

Pakiet ma dwie odpowiedzialności. Tor do hosta zdalnego przekłada polecenie procesu okna na wywołanie
SSH, którym proces rusza na maszynie wskazanej przez operatora, oraz przenosi pliki tym samym torem.
Dosięgnięcie operatora obejmuje kolejkę powiadomień, rejestrację urządzeń do wołania, ponowienie
i wygaśnięcie; silnik doręcza po łączu, które operator już otworzył, i nie budzi wygaszonego telefonu,
ponieważ wymagałoby to usługi wypychania powiadomień spoza tego systemu.

SSH jest tu transportem, nie drugim wykonawcą. Pakiet nie uruchamia procesów: buduje wyłącznie wiersz
poleceń transportu, a startuje go jedyny spawner platformy. Po stronie zdalnej proces uruchamia usługę
SSH, którą host już wystawia, więc w drzewie nie przybywa żaden własny demon ani protokół. Strumienie
SSH są strumieniami procesu zdalnego: wejście, wyjście, wyjście diagnostyczne i kod zakończenia
przechodzą wprost, a zerwanie połączenia kończy proces po stronie zdalnej sygnałem zawieszenia. Rola
agenta produktu komponuje się z tym torem bez zmian: wywołanie roli agenta przez SSH wykonuje żądania
kontraktu na hoście zdalnym tym samym binarium.

Trzy granice ograniczają ten pakiet. Identyfikator uchwytu procesu jest identyfikatorem lokalnego
procesu transportu, nie procesu zdalnego; drzewo potomstwa obejmowane przez warstwę sesji kończy
transport, a proces zdalny kończy się z usługą SSH po zerwaniu połączenia. Dziedziczenie środowiska
rdzenia dotyczy maszyny rdzenia: na hoście zdalnym proces dziedziczy środowisko logowania SSH, a wpisy
własne polecenia jadą w komendzie zdalnej. Zgoda na hosta jest wierszem tabeli hostów zdalnych,
wydawanym przez operatora — domyślnie jej nie ma, a rdzeń nie zainicjuje połączenia, którego operator
nie oddał.

Pakiet czyta bazę rdzenia przez uchwyt podany funkcją Zasil, a każda droga wywołana przed zasileniem
odmawia, nazywając brakujące wpięcie. Uchwyt podaje kompozycja programu głównego, więc odmowa braku
zasilenia dotyczy wołających spoza kompozycji i sprawdzianów. Nadajnik powiadomień jest osobnym
wpięciem: bez niego przebieg kolejki odmawia w całości i nie tyka ani jednego wiersza, a powiadomienia
czekają z pełnym budżetem prób, zamiast po cichu wygasać.

## budowa/server/internal/dane/studio_wsad.go
Typ i kontrakt obszaru wsadu deklaruje plik studio.go; ten plik implementuje
wyłącznie metody obszaru wsadu na tym samym uchwycie repozytorium studia, tak
jak plik studio_wersje.go implementuje metody obszaru wersji. Przebieg wsadu
i jego pozycje zapisują się w jednej transakcji, ponieważ liczby przyjętych
i odrzuconych pozycji w nagłówku przebiegu są sumą wierszy pozycji: zapis
rozdzielony na dwa osobne polecenia zostawiałby nagłówek niezgodny z jego
własnymi pozycjami.

## budowa/shared/swiezosc_generatu_test.go

Plik kontraktu jest źródłem prawdy, ale w kompilacji nie bierze udziału: bierze w niej udział
wygenerowany kod rdzenia, a w budowaniu klienta wygenerowany kod klienta. Zmiana źródła bez puszczenia
generatora rozjeżdża je po cichu i obie strony kompilują się dalej — rdzeń zna nazwę, której klient nie
zna, albo odwrotnie; rozjazd wychodzi dopiero na gnieździe, u operatora maszyny. Sprawdzian puszcza
generator na kopii, a nie na katalogu źródłowym, ponieważ generator zapisuje artefakty na dysk, więc
puszczony na miejscu nadpisałby pliki sprawdzanego drzewa.

## budowa/server/internal/dane/studio_wersje.go
Typ i interfejs obszaru wersji deklaruje plik studio.go; ten plik implementuje
wyłącznie metody obszaru wersji na tym samym uchwycie repozytorium studia,
podobnie jak inne pliki obszarowe modułów danych implementują swój obszar na
współdzielonym repozytorium modułu.

Treść wersji niesie dwa pola: pole treści krótkiej wprost oraz pole odwołania
do pliku dla treści obszernej, tym samym sposobem co repozytorium wiadomości.

Przywrócenie wersji jest zapisem dwutabelowym: metoda PrzywrocWersje czyta
wersję docelową i nadpisuje treść dokumentu w jednej transakcji — bez niej
odczyt wersji i zapis dokumentu mogłyby rozjechać się przy równoległym zapisie
tego samego dokumentu z innego okna.

## budowa/server/internal/injection/zaczepy.go
Plik odbiera zdarzenia zaczepów z linii typu system strumienia programu claude
uruchomionego z przełącznikiem --include-hook-events. Koperta startu zaczepu
niesie pola subtype, hook_id, hook_name i hook_event; koperta odpowiedzi niesie
dodatkowo output, stdout, stderr, exit_code i outcome. Zdarzenie zaczepu jest
zdarzeniem wykonawczym opisującym przebieg działania zaczepu, nie treść
wypowiedzi modelu. Warstwa składania odpowiedzi otrzymuje zdarzenie haczykiem
NaZdarzenieZaczepu i prowadzi na jego podstawie dziennik zdarzeń oraz
diagnostykę; brak podpiętego haczyka nie zmienia przebiegu tury.

## budowa/server/internal/injection/zestaw_narzedzi.go
Rdzeń nie zna nazw narzędzi eksperta. Podzbiór wskazany definicją eksperta
rozstrzyga serwer narzędzi zapisany w pliku narzedzia/ekspert_wykaz.go,
ponieważ tylko on trzyma wykaz kontraktu i umie przełożyć kod eksperta na
pozycję albo na całą grupę. Gdyby rdzeń wyliczał tę listę samodzielnie,
powstałby drugi wykaz obok kontraktowego, dlatego stąd wychodzi wyłącznie
opis podstawy i dołożeń, a rozwinięcie opisu w listę należy do strony
przeciwnej. Pakiet składa wiersze uruchomienia procesów tury i nie zna ani
rdzenia, ani bazy, ani kontraktu, podobnie jak plik nakladka.go składa
warstwy promptu bez znajomości jego treści.

Nazwa przełącznika PrzelacznikDolozen jest umową dwóch pakietów: odczyt po
drugiej stronie w module cmd/danaco-narzedzia oraz w funkcji
narzedzia.RozbijDolozenia bierze tę nazwę stąd, zamiast zapisywać ją
osobno. Wiersz uruchomienia serwera narzędzi czyta flag.Parse na domyślnym
flag.CommandLine z trybem ExitOnError, więc nieznany przełącznik nie zostaje
pominięty, tylko kończy proces serwera narzędzi i pozbawia turę całego
wykazu narzędzi. Rozjazd dwóch zapisów tej samej nazwy jest więc awarią,
nie usterką kosmetyczną.

Pole ZestawTury niesie dwa źródła różnej natury: podstawa jest zawężeniem
wykazu przez definicję eksperta, a dołożenia są dokładaniem pozycji przez
operatora na czas sesji. Jedno pole nie wyraziłoby obu stanów, ponieważ
lista pusta znaczyłaby jednocześnie brak zawężenia i brak narzędzi, a te
dwa stany prowadzą do odmiennych skutków dla dostępności wykazu w turze.

Kolejność dołożeń w ZlozZestawTury zostaje kolejnością dokładania, zgodnie
z regułą stosowaną w pliku narzedzia/wykaz.go przy rozszerzeniu roli okna:
dołożenie dokłada pozycję, nie przestawia kolejności istniejących. Nazwa
powtórzona zostaje na pierwszej pozycji, na której się pojawiła; dołożenie
powtórzone jest czynnością pustą, nie błędem.

Kształt wyniku ArgumentyDolozen odpowiada funkcjom
narzedzia.ArgumentyZasiegu i ArgumentyEksperta, ponieważ wszystkie trzy
zasilają ten sam wpis danaco tą samą drogą argumentów uruchomienia.

## budowa/server/internal/konfig/definicje_aplikacji.go
Nastawy wykonania oraz izolacji rozstrzygają, jak platforma prowadzi
rozmowę; nastawy tego pliku rozstrzygają, jak stoi sam rdzeń, na przykład
czy nasłuch wymaga logowania. Poziom zasięgu, na którym mieszkają, jest
najszerszy i nie ma bytu nadrzędnego, ponieważ programu nie ma czym zawęzić.

Nastawy konta nadawczego wchodzą także zmiennymi środowiska przy starcie,
zapisanymi w pliku konfiguracja/srodowisko.go; zapis w tabeli ustawienie
przesłania środowisko, ponieważ pomyłki w adresie serwera poczty nie da się
naprawić bez zatrzymania rdzenia, gdyby jedyną drogą poprawki było
środowisko.

Wartość domyślna wymogu logowania jest pusta, nie fałsz, ponieważ nastawa
ma trzy stany: brak wskazania rozstrzyga adres nasłuchu, a wskazania „tak”
i „nie” ustalają wymóg wprost. Wartość domyślna fałsz skasowałaby stan
pierwszy i zniosłaby wymóg logowania na nasłuchu wystawionym poza pętlę
zwrotną.

## budowa/server/internal/konfig/katalog_definicji.go
Pakiet nie sięga do bazy samodzielnie: bierze pozycje katalogu przez
interfejs zrodloKatalogu, więc połączenie z bazą powstaje w punkcie
kompozycji rdzenia. Pozycje mają postać struktury SettingDefinition
kontraktu, dzięki czemu ta sama treść zasila zarówno rejestr rozstrzygania,
jak i okno konfiguracji, bez drugiego opisu tego samego ustawienia.

Katalog pusty albo niedostępny daje rejestr wbudowany rdzenia, ponieważ
brak katalogu jest brakiem pozycji do pokazania, a nie brakiem możliwości
rozstrzygnięcia nastawy.

## budowa/server/internal/konfig/kontekst_zasiegu.go
Porządek adresów zwracanych dla kontekstu jest dwupoziomowy: najpierw
poziom zasięgu, a w ramach poziomu oś, w kolejności konto, model,
platforma.

## budowa/server/internal/konfig/odwzorowanie_kontraktu.go
Poziom pusty w polu scope wpisu kontraktu nie jest błędem ani stanem
wyjątkowym; jest informacją, że wartość pochodzi z warstwy definicji, nie
z zapisanego ustawienia.

## budowa/server/internal/konfig/osie.go
Klucz rozstrzygania ustawienia jest złożony: klucz, poziom, byt poziomu, oś
i byt osi. Poziom rozstrzyga pierwszeństwo zawsze przed osią, więc
ustawienie zapisane per konto na poziomie globalnym nie bije ustawienia
zapisanego na oknie komunikacji. Zapis na oknie jest aktem najwęższym
i najbardziej celowym, a oś opisuje adresata wartości, nie jej wagę.
Odwrotna kolejność oznaczałaby, że wybór konta unieważnia decyzję podjętą
wprost w oknie, co odbierałoby operatorowi sterowanie zamiast je
rozszerzać.

## budowa/server/internal/konfig/poziomy.go
Żadna ścieżka pakietu nie odmawia rozstrzygnięcia; także błąd źródła danych
kończy się polityką domyślną.

Kolejność poziomów wyliczana przez poziomyOdNajwezszego jest odwrotnością
kolumny poziom_zasiegu.pierwszenstwo. Ta kolumna mieszka w bazie i niesie
własną numerację poziomów, od aplikacji przez wartość 0, globalnego przez
wartość 1, aż po okno przez wartość 8; powielanie tej numeracji w kodzie
byłoby drugą prawdą o tej samej kolejności.

## budowa/server/internal/konfig/rejestr_definicji.go
Objaśnienie jest częścią definicji, nie dodatkiem: odpowiada kolumnie
objasnienie ustawionej jako NOT NULL w modelu danych. Pusty zbiór
dopuszczalnych adresów przy definicji zbudowanej w kodzie nie jest bramą
zamykającą dostęp, tylko brakiem wskazania miejsca w oknie konfiguracji.

## budowa/server/internal/konfig/rozgloszenie.go
Rozgłośnia nie jest drugim mechanizmem nastaw ani drugą tabelą: nie
przechowuje ani jednej wartości ustawienia z własnej woli, tylko przy
każdym ogłoszeniu pyta ten sam rozstrzygacz o rozstrzygnięcie w kontekście
nasłuchującego i podaje dalej wynik. Jedynym stanem, jaki trzyma, jest
zapis tego, co już powiedziała, po to, by nie budzić nasłuchującego zmianą,
której nie było; ten zapis nie jest źródłem wartości, jest pamięcią
rozmowy.

Rozgłośnia nie jest też własną usługą ani własnym wątkiem: nie odpala
gorutyny, nie odpytuje niczego w pętli i nie ma zegara. Doręczenie dzieje
się w wątku tego, kto ogłosił zapis, czyli w torze komendy config.set.
Droga zapisu woła Oglos z kluczem, który się zmienił, dopiero po udanym
utrwaleniu wiersza tabeli ustawienie, ponieważ nastawa, która nie usiadła
w bazie, nie jest zmianą nastawy. Punkty wywołania leżą poza tym pakietem.

Pierwsze doręczenie w funkcji Sledz jest częścią umowy, nie uprzejmością.
Nasłuchujący, który nie ma się przeładowywać, musi skądś wziąć punkt
wyjścia; gdyby brał go osobnym pytaniem, miałby dwie drogi do jednej
wartości i wyścig między nimi, w którym zapis mieszczący się pomiędzy
pytaniem a zapisaniem się zginąłby. Jedna droga niesie jedno źródło.

Funkcja Oglos mówi, że wskazane klucze mogły się zmienić, a nie że się
zmieniły, ponieważ ogłaszający zna adres zapisu, a nie skutek dla każdego
nasłuchującego z osobna. Skutek liczy rozstrzygacz osobno dla kontekstu
każdego nasłuchu, dzięki czemu zapis na poziomie okna nie budzi nasłuchu
poziomu aplikacji, a zapis na poziomie aplikacji nie budzi nikogo, kto ma
wartość z węższego poziomu.

Uchwyt Nastawa istnieje obok funkcji Rozstrzygnij, ponieważ korzystający
z niego siedzi na drodze gorącej, gdzie straż bramki rozstrzyga przy
każdym pakiecie z gniazda, a Rozstrzygnij za każdym razem schodzi po
wartość do warstwy trwałości. Uchwyt zdejmuje ten koszt bez zdejmowania
prawdy: nie jest drugą wartością mogącą rozjechać się ze źródłem, ponieważ
jedyną drogą jego zmiany jest ogłoszenie ze źródła, a własnego zapisu
uchwyt nie przyjmuje.

## budowa/server/internal/konfig/rozstrzyganie.go
Rozstrzygnij i rozgłośnia liczą wartość tym samym rozstrzyganiem z tego
samego źródła, więc drugiej prawdy o nastawie nie ma. Wartość zerowa pola
rozgloszenia jest zdatna do pracy, więc rozstrzygacz zbudowany bez
nasłuchów niczego nie kosztuje.

Błąd zwracany przez funkcję zapisy wraca obok wyniku rozstrzygnięcia
i służy wyłącznie diagnostyce; nie zatrzymuje samego rozstrzygania.

## budowa/server/internal/konfig/wartosc.go
Funkcja KodujJSON nigdy nie zawodzi: wartość uszkodzona trafia do koperty
kontraktu jako napis, nie jako błąd.

## budowa/server/internal/konfig/zrodlo_pamieciowe.go
Rozstrzygacz zbudowany bez źródła sięga po puste źródło pamięciowe, dzięki
czemu brak warstwy trwałości nie blokuje startu rdzenia, tylko daje
politykę domyślną.

## budowa/server/internal/konfig/zrodlo_ustawien.go
Implementacja interfejsu Zrodlo czytająca tabelę ustawienie należy do
warstwy repozytoriów; implementacja pamięciowa z tego pakietu obsługuje
pracę bez trwałości i bez sprawdzenia.

Wpisy zwrócone mimo błędu źródła wchodzą do rozstrzygnięcia, a ustawienia
bez zapisu schodzą na wartości domyślne. Błąd jest widoczny w polityce
efektywnej jako informacja diagnostyczna, nie jako odmowa, więc
implementacja może zwrócić wynik częściowy razem z błędem.

## budowa/server/internal/dane/design_prototyp.go
Ramki wiąże się identyfikatorem zewnętrznym, a nie kluczem obcym: tymi samymi
wartościami operuje kontrakt i tymi samymi wraca `design.prototype.get`, więc
przekład klucza w obie strony nie miałby odbiorcy. Spójność pilnuje adapter —
obie ramki połączenia muszą leżeć w tej samej kompozycji.

Ramek nieosiągalnych repozytorium nie liczy, ponieważ osierocenie ramki jest
wnioskiem chwilowym z odczytu grafu, a nie trwałą cechą wiersza — ramka
osierocona dziś bywa jutro ramką początkową, więc utrwalanie tej cechy
w kolumnie oznaczałoby przechowywanie wniosku, który starzeje się bez zapisu.

## budowa/server/internal/konfiguracja/argumenty.go
Przełącznik wymogu logowania nie mógł powstać jako flaga logiczna, ponieważ
flag.Bool umiałby wyrazić tylko dwa stany i zamieniłby brak wskazania we
wskazanie „nie", zdejmując wymóg logowania wystawionemu rdzeniowi przez
samo pominięcie przełącznika w wywołaniu.

Funkcja wykazPoPrzecinku odrzuca człony puste, ponieważ wzorzec pusty
pasowałby do niczego, a w bibliotece gniazda do czegokolwiek, co byłoby
zachowaniem sprzecznym z intencją filtra pochodzenia.

## budowa/server/internal/dane/design_szablony.go
Szablon i prompt wydany mają ten sam kształt kontraktu `DesignPrompt`, ale różne
życie: szablon Operator nadpisuje, a prompt wydany jest zapisem tego, co się
stało. Powód rozdziału tabel stoi w nagłówku migracji 231.

Historia niesie prowenancję: prompt wydany wraca wraz z kodami zasobów, które
z niego powstały. Bez tego pole `DesignAsset.PromptId` byłoby kodem bez drugiej
strony — kontrakt komendy `design.prompt.history.list` ten przekład wnosi.

## budowa/server/internal/dane/design_wersje.go
Wersja jest migawką układu, nie odwołaniem do warstw żywych. Warstwy zapisuje
się kolumna w kolumnę, bo `design.board.update` usuwa je i wstawia od nowa przy
każdym zapisie — odwołanie wskazywałoby wtedy wiersze, których już nie ma,
a wersja przestałaby opisywać cokolwiek dokładnie wtedy, gdy jest potrzebna.

Wykaz wersji warstw nie czyta (kontrakt: `versions` bez `layers`), stąd
`liczba_warstw` utrwalona w wierszu wersji. Warstwy wchodzą wyłącznie przy
przywróceniu, osobnym odczytem `WarstwyWersjiKompozycjiDesignu`.

Wersja jest zawsze nowa. Nadpisania nie ma i nie ma być: wersja to zapis stanu
z konkretnej chwili, a nadpisanie oznaczałoby, że stan sprzed godziny właśnie
się zmienił.

## budowa/server/internal/konfiguracja/katalog_klienta.go
Nazwa katalogu klient jest nazwą katalogu projektu interfejsu w drzewie
budowy, więc rdzeń uruchomiony z korzenia tego drzewa znajduje pakiet bez
przełącznika. Ta sama nazwa wiąże wdrożenie: pakowanie wydania kładzie
pakiet obok binarium pod tą nazwą, nie pod własną, ponieważ dwie różne
nazwy tego samego katalogu wracałyby odmową pakietu przy pierwszym
uruchomieniu. Katalog nieistniejący nie wstrzymuje startu, ponieważ gniazdo
pracuje wtedy bez plików statycznych.

## budowa/server/internal/dane/studio_postac_malarz.go
Migracja 368 zapisała powód wprost, a ten plik go wykonuje: malarz kopiuje
postać, nie treść, i nanosi ją w innym miejscu, więc między pobraniem
a naniesieniem stoją dwie osobne komendy, `studio.format.painter.copy`
i `studio.format.painter.apply`. Postać trzymana w pamięci procesu przepadała
przy przeładowaniu rdzenia — Operator pobierał postać, rdzeń wstawał od nowa,
a naniesienie odmawiało, nie znajdując takiej postaci. Przy pracy modelu
przepadała jeszcze łatwiej: model pobiera postać jednym narzędziem i nanosi
drugim, być może po kilku innych czynnościach.

Malarz jest narzędziem jednej czynności. Postać pobrana wczoraj i naniesiona
dziś byłaby zaskoczeniem, nie pomocą — stąd kolumna `wygasa` i odczyt, który
wpisu wygasłego nie oddaje. Wygasły wiersz nie jest przy tym kasowany
w odczycie: sprzątanie idzie osobnym wywołaniem, bo odczyt, który po cichu
usuwa wiersze, jest odczytem zmieniającym stan.

Malarz przenosi postać między dokumentami — to jego zwykłe użycie w pakiecie
biurowym. Kluczem jest więc okno, w którym Operator pracuje; dokument,
z którego postać zabrano, stoi obok jako wiedza, a nie jako warunek.

Nazwy pomocnicze tego pliku niosą przedrostek `malarz`, ponieważ przestrzeń
nazw pakietu `dane` jest dzielona z innymi wykonawcami.

## budowa/server/internal/konfiguracja/srodowisko.go
Odczyt zmiennej środowiska z pominięciem wykazu zwracanego przez
ZmienneSrodowiska rozjeżdża wzorzec .env.example z implementacją. Pilnuje
tego granica pakietu: nazwy zmiennych są nieeksportowane, więc odczyt po
nazwie dosłownej spoza tego pakietu jest widoczny w przeglądzie zmian.

Wartość nieczytelna zmiennej startTLS zatrzymuje start rdzenia zamiast po
cichu znaczyć „nie”, ponieważ ciche zejście do rozmowy otwartym tekstem
oddałoby poświadczenie nadawcy każdemu po drodze.

Wartość nieczytelna zmiennej wystawienia na wszystkie interfejsy zatrzymuje
start z tego samego powodu: pomyłka w zapisie tej jednej zmiennej
rozstrzyga o tym, czy rdzeń stanie w sieci, czy na pętli zwrotnej, a
milczące „nie” byłoby tu najgorszym z możliwych wyników, w stronę
przeciwną niż przy TLS.

Wartość nieczytelna zmiennej wymogu logowania zatrzymuje start z tego
samego powodu: literówka rozstrzyga o tym, czy rdzeń pyta wołającego
o token, czy nie pyta nikogo o nic.
