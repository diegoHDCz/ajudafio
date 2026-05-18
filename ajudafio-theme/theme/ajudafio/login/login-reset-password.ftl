<#import "template.ftl" as layout>
<@layout.registrationLayout displayInfo=false displayWide=false; section>
    <#if section = "form">
        <div class="af-logo-wrap">
            <img src="${url.resourcesPath}/img/logo.png" alt="Ajuda Fio" class="af-form-logo" />
        </div>
        <div class="af-reset-icon">📧</div>
        <h1 class="af-form-title">Recuperar senha</h1>
        <p class="af-form-sub">Informe seu e-mail e enviaremos um link para redefinir sua senha</p>

        <form id="kc-reset-password-form" action="${url.loginAction}" method="post">
            <div class="af-field">
                <label for="username">
                    <#if !realm.loginWithEmailAllowed>Usuário<#elseif !realm.registrationEmailAsUsername>Usuário ou e-mail<#else>E-mail</#if>
                </label>
                <input type="text" id="username" name="username" autofocus
                    value="${(auth.attemptedUsername!'')}"
                    placeholder="<#if !realm.loginWithEmailAllowed>seu usuário<#elseif !realm.registrationEmailAsUsername>usuário ou e-mail<#else>seu@email.com</#if>" />
            </div>

            <button type="submit" class="af-btn">Enviar e-mail de recuperação</button>
        </form>

        <div class="af-register-link">
            <span>Lembrou sua senha?</span>
            <a href="${url.loginUrl}">Voltar ao login</a>
        </div>
    </#if>
</@layout.registrationLayout>
