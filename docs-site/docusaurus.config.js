// @ts-check
// Note: type annotations allow type checking and IDEs autocompletion

const lightCodeTheme = require('prism-react-renderer').themes.github;
const darkCodeTheme = require('prism-react-renderer').themes.dracula;

/** @type {import('@docusaurus/types').Config} */
const config = {
  title: 'Nexus Protocol',
  tagline: '🚀 Enterprise Application Protocol для AI-платформ и интеграций',
  favicon: 'img/favicon.ico',

  // Set the production url of your site here
  url: 'https://nexus-protocol.dev',
  // Set the /<baseUrl>/ pathname under which your site is served
  baseUrl: '/',

  // GitHub pages deployment config.
  organizationName: 'nexus-protocol',
  projectName: 'nexus-protocol-docs',

  onBrokenLinks: 'throw',
  onBrokenMarkdownLinks: 'warn',

  // Even if you don't use internalization, you can use this field to set useful
  // metadata like html lang. For example, if your site is Chinese, you may want
  // to replace "en" with "zh-Hans".
  i18n: {
    defaultLocale: 'ru',
    locales: ['ru'],
  },

  presets: [
    [
      'classic',
      /** @type {import('@docusaurus/preset-classic').Options} */
      ({
        docs: {
          sidebarPath: require.resolve('./sidebars.js'),
          // Please change this to your repo.
          editUrl: 'https://github.com/nexus-protocol/docs/tree/main/',
          routeBasePath: '/',
        },
        blog: false,
        theme: {
          customCss: require.resolve('./src/css/custom.css'),
        },
      }),
    ],
  ],

  themeConfig:
    /** @type {import('@docusaurus/preset-classic').ThemeConfig} */
    ({
      // Replace with your project's social card
      image: 'img/nexus-social-card.jpg',
      navbar: {
        title: 'Nexus Protocol',
        // Логотип убран, используется только текст с иконкой из Lucide (Network)
        items: [
          {
            type: 'doc',
            docId: 'index',
            position: 'left',
            label: 'Главная',
          },
          {
            href: 'https://github.com/nexus-protocol',
            label: 'GitHub',
            position: 'right',
          },
        ],
      },
      footer: {
        style: 'dark',
        links: [
          {
            title: 'Документация',
            items: [
              {
                label: 'Главная',
                to: '/',
              },
            ],
          },
          {
            title: 'Ресурсы',
            items: [
              {
                label: 'GitHub',
                href: 'https://github.com/nexus-protocol',
              },
              {
                label: 'API Reference',
                to: '/api-reference',
              },
            ],
          },
          {
            title: 'Поддержка',
            items: [
              {
                label: 'Email',
                href: 'mailto:contact@nexus.dev',
              },
              {
                label: 'Website',
                href: 'https://nexus.dev',
              },
            ],
          },
        ],
        copyright: `Copyright © ${new Date().getFullYear()} Nexus Protocol. Built with Docusaurus.`,
      },
      prism: {
        theme: lightCodeTheme,
        darkTheme: darkCodeTheme,
        additionalLanguages: ['go', 'bash', 'json', 'yaml'],
      },
    }),
};

module.exports = config;

