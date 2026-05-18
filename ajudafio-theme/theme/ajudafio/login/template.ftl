<#macro registrationLayout bodyClass="" displayInfo=false displayMessage=true displayRequiredFields=false displayWide=false showAnotherWayIfPresent=true socialProvidersNode="" loginPage="" auth="" scripts=[]>
<!DOCTYPE html>
<html lang="pt-BR">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta name="robots" content="noindex, nofollow">
    <title><#if realm?? && realm.displayName??>${realm.displayName}<#else>Ajuda Fio</#if></title>
    <link rel="icon" href="${url.resourcesPath}/img/logo.png" type="image/png">
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
    <link href="https://fonts.googleapis.com/css2?family=Nunito:wght@400;600;700;800&family=Nunito+Sans:wght@400;600&display=swap" rel="stylesheet">
    <link rel="stylesheet" href="${url.resourcesPath}/css/login.css">
</head>
<body class="login-pf">

    <div class="af-bg-shapes">
        <div class="af-shape af-shape-1"></div>
        <div class="af-shape af-shape-2"></div>
        <div class="af-shape af-shape-3"></div>
        <div class="af-shape af-shape-4"></div>
        <div class="af-dots af-dots-1"></div>
        <div class="af-dots af-dots-2"></div>
    </div>

    <div class="af-layout">

        <div class="af-side">
            <div class="af-side-content">
                <img src="${url.resourcesPath}/img/logo.png" alt="Ajuda Fio" class="af-side-logo" />
                <h2 class="af-side-title">Cuidando de quem<br>você <em>ama</em></h2>
                <p class="af-side-sub">Conecte-se à plataforma de cuidados da Ajuda Fio e gerencie tudo com segurança e carinho.</p>
                <div class="af-side-badges">
                    <div class="af-badge"><span>🔒</span> Seguro</div>
                    <div class="af-badge"><span>💙</span> Confiável</div>
                    <div class="af-badge"><span>🌿</span> Cuidado</div>
                </div>
            </div>
        </div>

        <div class="af-main">
            <div class="af-card">

                <#if displayMessage && message?? && message?has_content>
                    <#if !message.type?? || message.type != 'warning' || !isAppInitiatedAction??>
                        <div class="alert alert-${message.type!'info'}">
                            <span><#if message.summary??>${message.summary}</#if></span>
                        </div>
                    </#if>
                </#if>

                <#nested "form">

                <#if displayInfo>
                    <div class="af-info">
                        <#nested "info">
                    </div>
                </#if>

            </div>

            <div class="af-footer">
                &copy; ${.now?string('yyyy')} Ajuda Fio &mdash; Todos os direitos reservados
            </div>
        </div>

    </div>

    <script src="${url.resourcesPath}/js/login.js"></script>
</body>
</html>
</#macro>
